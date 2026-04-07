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
  <div class="min-h-screen p-6 max-w-3xl mx-auto">
    <!-- Back -->
    <Button
      label="← Home"
      text
      size="small"
      class="mb-6"
      @click="router.push('/')"
    />

    <LobbyStatus />

    <!-- Error -->
    <Message v-if="error" severity="error" class="mt-4">
      {{ error }}
    </Message>

    <!-- Game Master view -->
    <div v-if="isGameMaster" class="mt-8">
      <!-- GM hasn't chosen a side yet -->
      <template v-if="gmPlayerRole === null">
        <Message severity="info">
          You are the Game Master. Choose your side, then start the game once
          the other player has joined.
        </Message>
        <RoleSelector @select="handleRoleSelect" />
      </template>

      <!-- GM has chosen a side -->
      <template v-else>
        <Message severity="success">
          You are the Game Master playing as
          <strong>{{ gmPlayerRole }}</strong>.
        </Message>
      </template>

      <Button
        label="Start Game"
        class="mt-6 w-full"
        size="large"
        :disabled="!canStart"
        @click="handleStart"
      />
    </div>

    <!-- Non-GM: no role yet -->
    <div v-else-if="role === null">
      <RoleSelector @select="handleRoleSelect" />
    </div>

    <!-- Non-GM: already has a role -->
    <div v-else class="mt-8">
      <Message severity="success">
        You've joined as
        <strong>{{ role }}</strong>. Waiting for the Game Master to
        start the game…
      </Message>
    </div>
  </div>
</template>
