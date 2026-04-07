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
