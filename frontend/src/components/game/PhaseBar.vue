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
