<script setup lang="ts">
import { onMounted } from 'vue';
import { usePlayerStore } from '@/stores/usePlayerStore';

const playerStore = usePlayerStore();

onMounted(async () => {
  await playerStore.initPlayer();
});
</script>

<template>
  <div class="dark">
    <RouterView />
  </div>
</template>
