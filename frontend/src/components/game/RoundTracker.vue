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
