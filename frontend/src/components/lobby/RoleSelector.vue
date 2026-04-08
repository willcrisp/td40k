<script setup lang="ts">
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
import { useRoomStore } from '@/stores/useRoomStore';
import { usePlayerStore } from '@/stores/usePlayerStore';

const emit = defineEmits<{
  (e: 'select', role: 'attacker' | 'defender'): void;
}>();

const room = useRoomStore();
const player = usePlayerStore();
const { attackerId, defenderId } = storeToRefs(room);
const { playerId } = storeToRefs(player);

function isMe(id: string | null) {
  return id === playerId.value;
}
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Attacker -->
    <div
      class="role-card"
      :class="{ 'role-card--selected': isMe(attackerId) }"
    >
      <div class="role-card-accent role-card-accent--attacker"></div>
      <div class="role-card-body">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-display text-primary" style="font-size: 1.125rem">
            Attacker
          </h3>
          <span
            v-if="attackerId"
            class="font-mono text-xs"
            :class="isMe(attackerId) ? 'text-primary' : 'text-surface-variant'"
          >
            {{ isMe(attackerId) ? 'YOU' : 'TAKEN' }}
          </span>
        </div>
        <p class="text-sm font-mono text-surface-variant mb-6">
          Controls the attacking force.
        </p>
        <Button
          v-if="!attackerId"
          label="Deploy as Attacker"
          class="btn-tactical w-full"
          @click="emit('select', 'attacker')"
        />
      </div>
    </div>

    <!-- Defender -->
    <div
      class="role-card"
      :class="{ 'role-card--selected': isMe(defenderId) }"
    >
      <div class="role-card-accent role-card-accent--defender"></div>
      <div class="role-card-body">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-display text-secondary" style="font-size: 1.125rem">
            Defender
          </h3>
          <span
            v-if="defenderId"
            class="font-mono text-xs"
            :class="isMe(defenderId) ? 'text-secondary' : 'text-surface-variant'"
          >
            {{ isMe(defenderId) ? 'YOU' : 'TAKEN' }}
          </span>
        </div>
        <p class="text-sm font-mono text-surface-variant mb-6">
          Controls the defending force.
        </p>
        <Button
          v-if="!defenderId"
          label="Deploy as Defender"
          class="btn-tactical w-full"
          style="background-color: var(--secondary-container) !important"
          @click="emit('select', 'defender')"
        />
      </div>
    </div>
  </div>
</template>
