<script setup lang="ts">
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
import Card from 'primevue/card';
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
  <div class="grid grid-cols-2 gap-6 mt-6">
    <!-- Attacker -->
    <Card
      :class="[
        'border-2 transition-colors',
        isMe(attackerId) ? 'border-red-500' : 'border-transparent',
      ]"
    >
      <template #title>
        <span class="text-red-400">⚔ Attacker</span>
      </template>
      <template #content>
        <p class="text-sm text-surface-400 mb-4">
          Controls the attacking force.
        </p>
        <div v-if="attackerId">
          <span class="text-green-400 text-sm">
            {{ isMe(attackerId) ? '✅ You' : '✅ Taken' }}
          </span>
        </div>
        <Button
          v-else
          label="Choose Attacker"
          severity="danger"
          @click="emit('select', 'attacker')"
        />
      </template>
    </Card>

    <!-- Defender -->
    <Card
      :class="[
        'border-2 transition-colors',
        isMe(defenderId) ? 'border-blue-500' : 'border-transparent',
      ]"
    >
      <template #title>
        <span class="text-blue-400">🛡 Defender</span>
      </template>
      <template #content>
        <p class="text-sm text-surface-400 mb-4">
          Controls the defending force.
        </p>
        <div v-if="defenderId">
          <span class="text-green-400 text-sm">
            {{ isMe(defenderId) ? '✅ You' : '✅ Taken' }}
          </span>
        </div>
        <Button
          v-else
          label="Choose Defender"
          severity="info"
          @click="emit('select', 'defender')"
        />
      </template>
    </Card>
  </div>
</template>
