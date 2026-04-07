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
import GameSidebar from '@/components/game/GameSidebar.vue';
import PhaseController from '@/components/game/PhaseController.vue';

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
  <div class="game-layout">
    <!-- Left Sidebar -->
    <GameSidebar />

    <!-- Main Content Area -->
    <div class="game-main">
      <!-- Top Header + Phase Bar -->
      <GameHUD />

      <!-- Canvas (front and center) -->
      <div class="game-canvas-area">
        <BoardCanvas />
      </div>

      <!-- Bottom Phase Controller -->
      <PhaseController />
    </div>
  </div>
</template>

<style scoped>
.game-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background-color: var(--surface-container-lowest);
}

.game-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background-color: var(--surface);
}

.game-canvas-area {
  flex: 1;
  position: relative;
  overflow: hidden;
  background-color: var(--surface-container-lowest);
  /* The canvas area should feel like it's recessed */
  box-shadow: inset 0 0 40px rgba(0, 0, 0, 0.5);
}
</style>
