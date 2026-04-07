<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useRoomStore } from '@/stores/useRoomStore';
import { usePlayerStore } from '@/stores/usePlayerStore';
import { useWebSocketStore } from '@/stores/useWebSocketStore';
import { useBoardStore } from '@/stores/useBoardStore';
import BoardCanvas from '@/components/game/BoardCanvas.vue';
import GameHUD from '@/components/game/GameHUD.vue';

const route = useRoute();
const roomStore = useRoomStore();
const playerStore = usePlayerStore();
const wsStore = useWebSocketStore();
const boardStore = useBoardStore();

const roomId = route.params.id as string;
const { playerId } = storeToRefs(playerStore);

onMounted(async () => {
  boardStore.reset();
  await roomStore.loadRoom(roomId);
  wsStore.connect(roomId, playerId.value);
});

onUnmounted(() => {
  wsStore.disconnect();
});
</script>

<template>
  <div class="relative w-screen h-screen overflow-hidden bg-surface-950">
    <BoardCanvas />
    <GameHUD />
  </div>
</template>
