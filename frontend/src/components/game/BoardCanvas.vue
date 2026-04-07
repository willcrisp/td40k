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
