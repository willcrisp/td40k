<script setup lang="ts">
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import { useRouter } from 'vue-router';
import { useWebSocketStore } from '@/stores/useWebSocketStore';
import { useRoomStore } from '@/stores/useRoomStore';
import PhaseBar from './PhaseBar.vue';
import RoundTracker from './RoundTracker.vue';
import PhaseController from './PhaseController.vue';

const router = useRouter();
const ws = useWebSocketStore();
const room = useRoomStore();
const { connected } = storeToRefs(ws);
const { name, status, winner } = storeToRefs(room);

function goHome() {
  ws.disconnect();
  router.push('/');
}
</script>

<template>
  <div
    class="absolute inset-x-0 top-0 flex flex-col gap-2 p-3
           bg-surface-900/90 backdrop-blur-sm border-b border-surface-700"
  >
    <!-- Top row -->
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <Button
          label="← Home"
          text
          size="small"
          @click="goHome"
        />
        <span class="text-sm font-bold truncate max-w-40">
          {{ name }}
        </span>
        <Tag
          v-if="!connected"
          value="Reconnecting…"
          severity="warn"
          class="text-xs"
        />
      </div>

      <RoundTracker />

      <PhaseController />
    </div>

    <!-- Phase bar -->
    <PhaseBar />

    <!-- Game over banner -->
    <div
      v-if="status === 'finished'"
      class="text-center py-1 bg-yellow-800/80 rounded text-sm font-bold"
    >
      ⚔ Game Over
      <span v-if="winner">— {{ winner.toUpperCase() }} wins</span>
    </div>
  </div>
</template>
