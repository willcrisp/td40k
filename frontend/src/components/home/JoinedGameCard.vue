<script setup lang="ts">
import type { JoinedGameSummary } from '@/types';
import { PHASE_LABELS } from '@/types';
import Button from 'primevue/button';

const props = defineProps<{ game: JoinedGameSummary }>();
const emit = defineEmits<{
  (e: 'rejoin', id: string): void;
}>();

const statusColorClass: Record<string, string> = {
  lobby: 'text-tertiary',
  active: 'text-primary',
  finished: 'text-secondary',
  closed: 'text-surface-variant',
};

const roleLabel: Record<string, string> = {
  attacker: 'ATTACKER',
  defender: 'DEFENDER',
};
</script>

<template>
  <div class="panel-low p-0 riveted flex flex-col h-full border-l-4 border-outline-variant">
    <div class="p-card-header-bar flex justify-between items-center bg-surface-container-highest">
      <h3 class="font-display text-lg text-white truncate">{{ props.game.name }}</h3>
      <span class="text-[10px] font-mono px-2 py-0.5 bg-secondary-container text-white">
        {{ roleLabel[props.game.role] }}
      </span>
    </div>
    
    <div class="p-4 flex-1">
      <div class="mb-4">
        <span 
          class="text-xs font-mono px-2 py-1 bg-surface-container-highest border border-ghost-border"
          :class="statusColorClass[props.game.status]"
        >
          {{ props.game.status.toUpperCase() }}
        </span>
      </div>

      <div class="font-mono text-sm space-y-2">
        <template v-if="props.game.status === 'active' || props.game.status === 'lobby'">
          <p class="text-tertiary">Round {{ props.game.battle_round }} of 5</p>
          <p class="text-surface-variant uppercase">{{ PHASE_LABELS[props.game.current_phase] }}</p>
        </template>
        
        <p v-else-if="props.game.status === 'finished'" class="text-tertiary">
          GAME OVER // FINISHED
        </p>
        <p v-else-if="props.game.status === 'closed'" class="text-surface-variant">
          GAME CLOSED // OFFLINE
        </p>
      </div>
    </div>

    <div class="p-4 bg-surface-container-highest mt-auto">
      <Button
        label="REJOIN"
        class="btn-tactical w-full"
        :disabled="props.game.status === 'closed' || props.game.status === 'finished'"
        @click="emit('rejoin', props.game.id)"
      />
    </div>
  </div>
</template>
