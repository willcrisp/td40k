<script setup lang="ts">
import { onMounted } from 'vue';
import { usePlayerStore } from '@/stores/usePlayerStore';

const playerStore = usePlayerStore();

onMounted(async () => {
  await playerStore.initPlayer();
  // Ensure dark mode is applied to document
  document.documentElement.classList.add('dark');
});
</script>

<template>
  <div class="dark">
    <RouterView />
  </div>
</template>
