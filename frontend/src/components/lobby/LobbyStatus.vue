<script setup lang="ts">
import { storeToRefs } from 'pinia';
import Tag from 'primevue/tag';
import { useRoomStore } from '@/stores/useRoomStore';

const room = useRoomStore();
const { attackerId, defenderId, name, status } = storeToRefs(room);
</script>

<template>
  <div class="flex flex-col gap-3">
    <h2 class="text-2xl font-bold">{{ name }}</h2>
    <div class="flex gap-4">
      <div class="flex items-center gap-2">
        <span class="text-sm text-surface-400">Attacker:</span>
        <Tag
          :value="attackerId ? 'Ready' : 'Waiting…'"
          :severity="attackerId ? 'success' : 'warn'"
        />
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm text-surface-400">Defender:</span>
        <Tag
          :value="defenderId ? 'Ready' : 'Waiting…'"
          :severity="defenderId ? 'success' : 'warn'"
        />
      </div>
    </div>
    <p v-if="status === 'lobby'" class="text-sm text-surface-400">
      Share this page URL for others to join.
    </p>
  </div>
</template>
