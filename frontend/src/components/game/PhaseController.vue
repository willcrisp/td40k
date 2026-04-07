<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { ref, computed } from 'vue';
import ConfirmDialog from 'primevue/confirmdialog';
import { useConfirm } from 'primevue/useconfirm';
import { useRoomStore } from '@/stores/useRoomStore';

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

  <footer
    v-if="isGameMaster"
    class="phase-controller panel-high border-t border-outline-variant"
  >
    <div class="flex items-center justify-between gap-6 w-full max-w-5xl px-8 flex-nowrap">
      <button
        class="btn-secondary-tactical px-6 h-12 flex items-center gap-3 cursor-pointer flex-shrink-0 whitespace-nowrap"
        :disabled="loading"
        @click="handlePrev"
      >
        <i class="pi pi-chevron-left text-sm" />
        <span class="font-display tracking-widest text-xs">PREVIOUS PHASE</span>
      </button>

      <div class="flex-1 flex flex-col items-center min-w-0">
        <p class="text-[10px] font-mono text-surface-variant uppercase whitespace-nowrap overflow-hidden text-ellipsis">
          GAME MASTER CONTROL
        </p>
        <div class="w-full h-[1px] bg-outline-variant opacity-20 mt-1"></div>
      </div>

      <button
        class="btn-tactical px-10 h-10 flex items-center gap-4 flex-shrink-0 whitespace-nowrap"
        :class="{ 'pulse-final': isFinalStep }"
        :disabled="loading"
        @click="handleNext"
      >
        <span class="font-display tracking-widest text-sm">NEXT PHASE</span>
        <i class="pi pi-chevron-right text-xs" />
      </button>
    </div>
  </footer>
</template>

<style scoped>
.phase-controller {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  min-height: 64px;
  z-index: 20;
}

.pulse-final {
  animation: pulse-danger 3s infinite steps(2);
}

@keyframes pulse-danger {
  0%, 100% { opacity: 1; filter: brightness(1); }
  50% { opacity: 0.6; filter: brightness(0.85); }
}
</style>
