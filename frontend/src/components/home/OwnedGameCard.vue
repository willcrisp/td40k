<script setup lang="ts">
import type { OwnedGameSummary } from '@/types';
import { PHASE_LABELS } from '@/types';
import Button from 'primevue/button';

const props = defineProps<{ game: OwnedGameSummary }>();
const emit = defineEmits<{
  (e: 'enter', id: string): void;
  (e: 'close', id: string): void;
}>();

const statusColorClass: Record<string, string> = {
  lobby: 'text-tertiary',
  active: 'text-primary',
  finished: 'text-secondary',
  closed: 'text-surface-variant',
};
</script>

<template>
  <div class="panel-low p-0 riveted flex flex-col h-full border-l-4 border-secondary-container">
    <div class="p-card-header-bar">
      <h3 class="font-display text-lg text-white truncate">{{ props.game.name }}</h3>
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
        <template v-if="props.game.status === 'active'">
          <p class="text-tertiary">Round {{ props.game.battle_round }} of 5</p>
          <p class="text-primary">{{ props.game.active_player.toUpperCase() }}'S TURN</p>
          <p class="text-surface-variant uppercase">{{ PHASE_LABELS[props.game.current_phase] }}</p>
        </template>
        
        <template v-else-if="props.game.status === 'lobby'">
          <p class="text-surface-variant">Waiting for players...</p>
          <div class="grid grid-cols-2 gap-2 mt-2">
            <div class="bg-surface-container-lowest p-2 border border-ghost-border">
              <p class="text-[10px] text-surface-variant">ATTACKER</p>
              <p class="text-xs" :class="props.game.attacker_id ? 'text-primary' : 'text-surface-variant opacity-50'">
                {{ props.game.attacker_id ? 'CONNECTED' : 'OPEN' }}
              </p>
            </div>
            <div class="bg-surface-container-lowest p-2 border border-ghost-border">
              <p class="text-[10px] text-surface-variant">DEFENDER</p>
              <p class="text-xs" :class="props.game.defender_id ? 'text-primary' : 'text-surface-variant opacity-50'">
                {{ props.game.defender_id ? 'CONNECTED' : 'OPEN' }}
              </p>
            </div>
          </div>
        </template>
        
        <p v-else-if="props.game.status === 'finished'" class="text-tertiary">
          GAME OVER // FINISHED
        </p>
        <p v-else-if="props.game.status === 'closed'" class="text-surface-variant">
          GAME CLOSED // OFFLINE
        </p>
      </div>
    </div>

    <div class="p-4 bg-surface-container-highest flex gap-3 mt-auto">
      <Button
        label="ENTER"
        class="btn-tactical flex-1"
        @click="emit('enter', props.game.id)"
      />
      <Button
        label="CLOSE"
        class="btn-secondary-tactical px-4"
        :disabled="props.game.status === 'closed' || props.game.status === 'finished'"
        @click="emit('close', props.game.id)"
      />
    </div>
  </div>
</template>
