<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
import Message from 'primevue/message';
import { useRoomStore } from '@/stores/useRoomStore';
import { usePlayerStore } from '@/stores/usePlayerStore';
import { useWebSocketStore } from '@/stores/useWebSocketStore';
import LobbyStatus from '@/components/lobby/LobbyStatus.vue';
import RoleSelector from '@/components/lobby/RoleSelector.vue';

const route = useRoute();
const router = useRouter();
const roomStore = useRoomStore();
const playerStore = usePlayerStore();
const wsStore = useWebSocketStore();

const roomId = route.params.id as string;
const { status, role, canStart, isGameMaster, gmPlayerRole } = storeToRefs(roomStore);
const { playerId } = storeToRefs(playerStore);

const error = ref<string | null>(null);

onMounted(async () => {
  await roomStore.loadRoom(roomId);

  // If already active and player is participant, go to game
  if (
    roomStore.status === 'active' &&
    (roomStore.attackerId === playerId.value ||
      roomStore.defenderId === playerId.value ||
      roomStore.gameMasterId === playerId.value)
  ) {
    router.replace(`/game/${roomId}`);
    return;
  }

  wsStore.connect(roomId, playerId.value);
});

onUnmounted(() => {
  wsStore.disconnect();
});

// Redirect when GM starts the game
watch(status, (val) => {
  if (val === 'active') {
    router.push(`/game/${roomId}`);
  }
});

async function handleRoleSelect(selectedRole: 'attacker' | 'defender') {
  error.value = null;
  try {
    await roomStore.joinRoom(roomId, selectedRole);
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? 'Failed to join';
  }
}

async function handleStart() {
  await roomStore.startGame();
}
</script>

<template>
  <div class="layout-terminal">
    <!-- Back to Hub -->
    <div class="mb-10">
      <Button
        label="&lt;&lt; BACK TO LOBBY"
        class="btn-secondary-tactical"
        @click="router.push('/')"
      />
    </div>

    <!-- Header Section -->
    <div class="mb-8 border-b border-outline-variant pb-4">
      <h1 class="text-3xl font-display text-primary">GAME LOBBY</h1>
      <p class="text-sm font-mono text-surface-variant">GAME ID: {{ roomId.split('-')[0].toUpperCase() }} // STATUS: CONNECTED</p>
    </div>

    <LobbyStatus />

    <!-- Error Messaging (System Alert) -->
    <div v-if="error" class="mt-6 panel-low p-4 border-l-4 border-primary riveted">
      <p class="text-xs font-mono text-primary mb-1">ERROR</p>
      <p class="text-sm font-mono">{{ error }}</p>
    </div>

    <!-- Interface Logic -->
    <div class="mt-10 max-w-2xl">
      <!-- Game Master view -->
      <div v-if="isGameMaster">
        <!-- GM Protocol Prompt -->
        <div v-if="gmPlayerRole === null" class="mb-8 panel-low p-6 border-l-4 border-tertiary riveted">
          <h2 class="text-xl font-display text-tertiary mb-2">SELECT YOUR ROLE</h2>
          <p class="text-sm text-surface-variant mb-6 uppercase">Choose a side to begin the game.</p>
          <RoleSelector @select="handleRoleSelect" />
        </div>

        <!-- GM Assigned State -->
        <div v-else class="mb-8 panel-lowest p-4 border-l-4 border-secondary-container">
          <p class="text-xs font-mono text-secondary mb-1">GAME MASTER</p>
          <p class="text-sm font-display text-white">YOUR ROLE: {{ gmPlayerRole.toUpperCase() }}</p>
        </div>

        <Button
          label="START GAME"
          class="btn-tactical w-full h-16 text-xl"
          :disabled="!canStart"
          @click="handleStart"
        />
        <p v-if="!canStart" class="text-[10px] font-mono text-center mt-2 text-surface-variant uppercase">
          Awaiting players to join before starting...
        </p>
      </div>

      <!-- Non-GM: no role yet -->
      <div v-else-if="role === null" class="panel-low p-6 border-l-4 border-tertiary riveted">
        <h2 class="text-xl font-display text-tertiary mb-2">SELECT YOUR ROLE</h2>
        <p class="text-sm text-surface-variant mb-6 uppercase">Choose a side to join the game.</p>
        <RoleSelector @select="handleRoleSelect" />
      </div>

      <!-- Non-GM: already has a role -->
      <div v-else class="panel-low p-6 border-l-4 border-secondary-container riveted">
        <p class="text-xs font-mono text-secondary mb-1">ROLE ASSIGNED</p>
        <h2 class="text-xl font-display text-white mb-2">YOUR ROLE: {{ role.toUpperCase() }}</h2>
        <p class="text-sm text-surface-variant uppercase">
          Waiting for the Game Master to start the game...
        </p>
      </div>
    </div>
  </div>
</template>
