<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useRoomStore } from '@/stores/useRoomStore';
import { PHASE_LABELS } from '@/types';

const room = useRoomStore();
const { battleRound, activePlayer, currentPhase } = storeToRefs(room);
</script>

<template>
  <div class="round-tracker">
    <button
      class="tracker-tab font-display"
    >
      ROUND {{ String(battleRound).padStart(2, '0') }}
    </button>

    <button
      class="tracker-tab font-display"
      :class="{ 'tracker-tab--attacker': activePlayer === 'attacker', 'tracker-tab--defender': activePlayer === 'defender' }"
    >
      {{ activePlayer === 'attacker' ? 'ATTACKER' : 'DEFENDER' }}
    </button>

    <button
      class="tracker-tab font-display tracker-tab--phase"
    >
      {{ PHASE_LABELS[currentPhase].toUpperCase() }}
    </button>
  </div>
</template>

<style scoped>
.round-tracker {
  display: flex;
  align-items: center;
  gap: 2px;
}

.tracker-tab {
  padding: 6px 16px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.02em;
  border: none;
  cursor: pointer;
  background: var(--color-surface-high);
  color: var(--color-on-surface-variant);
  transition: background 0.15s, color 0.15s;
  white-space: nowrap;
}

.tracker-tab:hover {
  background: var(--color-surface-highest);
  color: var(--color-on-surface);
}

.tracker-tab--attacker {
  color: var(--color-primary);
  background: rgba(163, 19, 23, 0.2);
}

.tracker-tab--defender {
  color: var(--color-secondary);
  background: rgba(10, 76, 106, 0.3);
}

.tracker-tab--phase {
  color: var(--color-on-surface);
  background: var(--color-surface-highest);
}
</style>
