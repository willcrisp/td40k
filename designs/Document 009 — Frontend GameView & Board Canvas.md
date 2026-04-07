Document 009 — Frontend: GameView & Board Canvas

Purpose


Implement the full-screen game view with the pan/zoom canvas board and the HUD overlay.


---

frontend/src/composables/useBoardControls.ts

	import { ref, type Ref } from 'vue';
	import { useBoardStore } from '@/stores/useBoardStore';
	
	const MIN_ZOOM = 0.2;
	const MAX_ZOOM = 8.0;
	
	export function useBoardControls(canvas: Ref<HTMLCanvasElement | null>) {
	  const board = useBoardStore();
	  const isDragging = ref(false);
	  let lastX = 0;
	  let lastY = 0;
	
	  function onWheel(e: WheelEvent) {
	    e.preventDefault();
	    if (!canvas.value) return;
	
	    const factor = e.deltaY < 0 ? 1.1 : 0.9;
	    const newZoom = Math.min(
	      MAX_ZOOM,
	      Math.max(MIN_ZOOM, board.zoom * factor)
	    );
	
	    const rect = canvas.value.getBoundingClientRect();
	    const mx = e.clientX - rect.left;
	    const my = e.clientY - rect.top;
	
	    // Zoom towards cursor
	    board.panX = mx - ((mx - board.panX) * newZoom) / board.zoom;
	    board.panY = my - ((my - board.panY) * newZoom) / board.zoom;
	    board.zoom = newZoom;
	  }
	
	  function onPointerDown(e: PointerEvent) {
	    isDragging.value = true;
	    lastX = e.clientX;
	    lastY = e.clientY;
	    canvas.value?.setPointerCapture(e.pointerId);
	  }
	
	  function onPointerMove(e: PointerEvent) {
	    if (!isDragging.value) return;
	    board.panX += e.clientX - lastX;
	    board.panY += e.clientY - lastY;
	    lastX = e.clientX;
	    lastY = e.clientY;
	  }
	
	  function onPointerUp(e: PointerEvent) {
	    isDragging.value = false;
	    canvas.value?.releasePointerCapture(e.pointerId);
	  }
	
	  return { onWheel, onPointerDown, onPointerMove, onPointerUp };
	}


---

frontend/src/components/game/BoardCanvas.vue

	<script setup lang="ts">
	import {
	  ref,
	  onMounted,
	  onUnmounted,
	  watchEffect,
	} from 'vue';
	import { useBoardStore } from '@/stores/useBoardStore';
	import { useBoardControls } from '@/composables/useBoardControls';
	
	const BOARD_W_IN = 44;
	const BOARD_H_IN = 60;
	const BASE_PX = 14; // pixels per inch at zoom 1.0
	
	const canvas = ref<HTMLCanvasElement | null>(null);
	const board = useBoardStore();
	const controls = useBoardControls(canvas);
	
	let animFrame: number;
	
	function resizeCanvas() {
	  if (!canvas.value) return;
	  canvas.value.width = canvas.value.offsetWidth;
	  canvas.value.height = canvas.value.offsetHeight;
	}
	
	function draw() {
	  const el = canvas.value;
	  if (!el) return;
	  const ctx = el.getContext('2d');
	  if (!ctx) return;
	
	  ctx.clearRect(0, 0, el.width, el.height);
	
	  ctx.save();
	  ctx.translate(board.panX, board.panY);
	  ctx.scale(board.zoom, board.zoom);
	
	  const w = BOARD_W_IN * BASE_PX;
	  const h = BOARD_H_IN * BASE_PX;
	
	  // Board background
	  ctx.fillStyle = '#1a3020';
	  ctx.fillRect(0, 0, w, h);
	
	  // Grid lines
	  ctx.strokeStyle = 'rgba(255,255,255,0.12)';
	  ctx.lineWidth = 0.5 / board.zoom;
	
	  for (let x = 0; x <= BOARD_W_IN; x++) {
	    ctx.beginPath();
	    ctx.moveTo(x * BASE_PX, 0);
	    ctx.lineTo(x * BASE_PX, h);
	    ctx.stroke();
	  }
	
	  for (let y = 0; y <= BOARD_H_IN; y++) {
	    ctx.beginPath();
	    ctx.moveTo(0, y * BASE_PX);
	    ctx.lineTo(w, y * BASE_PX);
	    ctx.stroke();
	  }
	
	  // Board border
	  ctx.strokeStyle = '#8b7355';
	  ctx.lineWidth = 3 / board.zoom;
	  ctx.strokeRect(0, 0, w, h);
	
	  // Inch labels at 6" intervals (optional, helps orientation)
	  ctx.fillStyle = 'rgba(255,255,255,0.3)';
	  ctx.font = `${10 / board.zoom}px monospace`;
	  for (let x = 0; x <= BOARD_W_IN; x += 6) {
	    ctx.fillText(`${x}"`, x * BASE_PX + 2, 10 / board.zoom);
	  }
	
	  ctx.restore();
	}
	
	function loop() {
	  draw();
	  animFrame = requestAnimationFrame(loop);
	}
	
	onMounted(() => {
	  resizeCanvas();
	  window.addEventListener('resize', resizeCanvas);
	  // Centre the board
	  if (canvas.value) {
	    board.panX = (canvas.value.width - BOARD_W_IN * BASE_PX) / 2;
	    board.panY = (canvas.value.height - BOARD_H_IN * BASE_PX) / 2;
	  }
	  loop();
	});
	
	onUnmounted(() => {
	  cancelAnimationFrame(animFrame);
	  window.removeEventListener('resize', resizeCanvas);
	});
	
	// Redraw whenever store changes (zoom/pan updates trigger a new frame anyway
	// via rAF loop, but watchEffect ensures reactivity is tracked)
	watchEffect(() => {
	  void board.zoom;
	  void board.panX;
	  void board.panY;
	});
	</script>
	
	<template>
	  <canvas
	    ref="canvas"
	    class="w-full h-full block cursor-grab active:cursor-grabbing"
	    @wheel="controls.onWheel"
	    @pointerdown="controls.onPointerDown"
	    @pointermove="controls.onPointerMove"
	    @pointerup="controls.onPointerUp"
	  />
	</template>


---

frontend/src/components/game/PhaseBar.vue

	<script setup lang="ts">
	import { storeToRefs } from 'pinia';
	import { PHASES, PHASE_LABELS, PHASE_NUMBERS } from '@/types';
	import { useRoomStore } from '@/stores/useRoomStore';
	
	const room = useRoomStore();
	const { currentPhase } = storeToRefs(room);
	</script>
	
	<template>
	  <div class="flex items-center gap-1">
	    <div
	      v-for="phase in PHASES"
	      :key="phase"
	      :class="[
	        'flex items-center gap-1 px-3 py-1 rounded text-xs font-bold uppercase',
	        'transition-colors duration-300',
	        currentPhase === phase
	          ? 'bg-red-600 text-white'
	          : 'bg-surface-800 text-surface-400',
	      ]"
	    >
	      <span>{{ PHASE_NUMBERS[phase] }}.</span>
	      <span class="hidden md:inline">
	        {{ PHASE_LABELS[phase].replace(' Phase', '') }}
	      </span>
	    </div>
	  </div>
	</template>


---

frontend/src/components/game/RoundTracker.vue

	<script setup lang="ts">
	import { storeToRefs } from 'pinia';
	import { useRoomStore } from '@/stores/useRoomStore';
	
	const room = useRoomStore();
	const { battleRound, activePlayer } = storeToRefs(room);
	
	const MAX_ROUNDS = 5;
	</script>
	
	<template>
	  <div class="flex items-center gap-4">
	    <!-- Round dots -->
	    <div class="flex items-center gap-1">
	      <span
	        v-for="r in MAX_ROUNDS"
	        :key="r"
	        :class="[
	          'w-3 h-3 rounded-full transition-colors',
	          r <= battleRound ? 'bg-red-500' : 'bg-surface-600',
	        ]"
	      />
	    </div>
	    <span class="text-sm font-semibold">
	      Round {{ battleRound }} of {{ MAX_ROUNDS }}
	    </span>
	    <span
	      :class="[
	        'text-sm font-bold uppercase px-2 py-0.5 rounded',
	        activePlayer === 'attacker' ? 'text-red-400' : 'text-blue-400',
	      ]"
	    >
	      {{ activePlayer }}'s Turn
	    </span>
	  </div>
	</template>


---

frontend/src/components/game/PhaseController.vue

	<script setup lang="ts">
	import { storeToRefs } from 'pinia';
	import { ref, computed } from 'vue';
	import Button from 'primevue/button';
	import ConfirmDialog from 'primevue/confirmdialog';
	import { useConfirm } from 'primevue/useconfirm';
	import { useRoomStore } from '@/stores/useRoomStore';
	import { PHASE_LABELS } from '@/types';
	
	const room = useRoomStore();
	const confirm = useConfirm();
	const { currentPhase, battleRound, activePlayer, isGameMaster } =
	  storeToRefs(room);
	
	const loading = ref(false);
	
	const isFinalStep = computed(
	  () =>
	    currentPhase.value === 'fight' &&
	    activePlayer.value === 'defender' &&
	    battleRound.value === 5
	);
	
	async function handleNext() {
	  if (isFinalStep.value) {
	    confirm.require({
	      message:
	        'This will end the game. ' +
	        'Are you sure you want to advance past the final Fight Phase?',
	      header: 'End Game?',
	      icon: 'pi pi-exclamation-triangle',
	      acceptLabel: 'End Game',
	      rejectLabel: 'Cancel',
	      accept: () => advance(),
	    });
	    return;
	  }
	  await advance();
	}
	
	async function advance() {
	  loading.value = true;
	  await room.nextPhase().finally(() => (loading.value = false));
	}
	
	async function handlePrev() {
	  loading.value = true;
	  await room.prevPhase().finally(() => (loading.value = false));
	}
	</script>
	
	<template>
	  <ConfirmDialog />
	
	  <div
	    v-if="isGameMaster"
	    class="flex items-center gap-3"
	  >
	    <Button
	      icon="pi pi-chevron-left"
	      severity="secondary"
	      rounded
	      :loading="loading"
	      @click="handlePrev"
	    />
	
	    <div class="text-center min-w-40">
	      <p class="text-xs text-surface-400 uppercase tracking-widest">
	        Current Phase
	      </p>
	      <p class="font-bold text-sm">
	        {{ PHASE_LABELS[currentPhase] }}
	      </p>
	    </div>
	
	    <Button
	      icon="pi pi-chevron-right"
	      :severity="isFinalStep ? 'danger' : 'primary'"
	      rounded
	      :loading="loading"
	      @click="handleNext"
	    />
	  </div>
	</template>


---

frontend/src/components/game/GameHUD.vue

	<script setup lang="ts">
	import { storeToRefs } from 'pinia';
	import Button from 'primevue/button';
	import Tag from 'primevue/tag';
	import { useRouter } from 'vue-router';
	import { useWebSocketStore } from '@/stores/useWebSocketStore';
	import { useRoomStore } from '@/stores/useRoomStore';
	import PhaseBar from './PhaseBar.vue';
	import RoundTracker from './RoundTracker.vue';
	import PhaseController from './PhaseController.vue';
	
	const router = useRouter();
	const ws = useWebSocketStore();
	const room = useRoomStore();
	const { connected } = storeToRefs(ws);
	const { name, status, winner } = storeToRefs(room);
	
	function goHome() {
	  ws.disconnect();
	  router.push('/');
	}
	</script>
	
	<template>
	  <div
	    class="absolute inset-x-0 top-0 flex flex-col gap-2 p-3
	           bg-surface-900/90 backdrop-blur-sm border-b border-surface-700"
	  >
	    <!-- Top row -->
	    <div class="flex items-center justify-between gap-4">
	      <div class="flex items-center gap-3">
	        <Button
	          label="← Home"
	          text
	          size="small"
	          @click="goHome"
	        />
	        <span class="text-sm font-bold truncate max-w-40">
	          {{ name }}
	        </span>
	        <Tag
	          v-if="!connected"
	          value="Reconnecting…"
	          severity="warn"
	          class="text-xs"
	        />
	      </div>
	
	      <RoundTracker />
	
	      <PhaseController />
	    </div>
	
	    <!-- Phase bar -->
	    <PhaseBar />
	
	    <!-- Game over banner -->
	    <div
	      v-if="status === 'finished'"
	      class="text-center py-1 bg-yellow-800/80 rounded text-sm font-bold"
	    >
	      ⚔ Game Over
	      <span v-if="winner">— {{ winner.toUpperCase() }} wins</span>
	    </div>
	  </div>
	</template>


---

frontend/src/views/GameView.vue

	<script setup lang="ts">
	import { onMounted, onUnmounted } from 'vue';
	import { useRoute } from 'vue-router';
	import { storeToRefs } from 'pinia';
	import { useRoomStore } from '@/stores/useRoomStore';
	import { usePlayerStore } from '@/stores/usePlayerStore';
	import { useWebSocketStore } from '@/stores/useWebSocketStore';
	import { useBoardStore } from '@/stores/useBoardStore';
	import BoardCanvas from '@/components/game/BoardCanvas.vue';
	import GameHUD from '@/components/game/GameHUD.vue';
	
	const route = useRoute();
	const roomStore = useRoomStore();
	const playerStore = usePlayerStore();
	const wsStore = useWebSocketStore();
	const boardStore = useBoardStore();
	
	const roomId = route.params.id as string;
	const { playerId } = storeToRefs(playerStore);
	
	onMounted(async () => {
	  boardStore.reset();
	  await roomStore.loadRoom(roomId);
	  wsStore.connect(roomId, playerId.value);
	});
	
	onUnmounted(() => {
	  wsStore.disconnect();
	});
	</script>
	
	<template>
	  <div class="relative w-screen h-screen overflow-hidden bg-surface-950">
	    <BoardCanvas />
	    <GameHUD />
	  </div>
	</template>