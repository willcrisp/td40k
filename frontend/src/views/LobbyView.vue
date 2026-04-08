<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
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
const { status, role, canStart, isGameMaster, gmPlayerRole, name } =
  storeToRefs(roomStore);
const { playerId, nickname, username } = storeToRefs(playerStore);

const error = ref<string | null>(null);
const attackerArmy = ref<string>('');
const defenderArmy = ref<string>('');

onMounted(async () => {
  await roomStore.loadRoom(roomId);

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

watch(status, (val) => {
  if (val === 'active') {
    router.push(`/game/${roomId}`);
  }
});

async function handleRoleSelect(
  selectedRole: 'attacker' | 'defender'
) {
  error.value = null;
  try {
    await roomStore.joinRoom(roomId, selectedRole);
  } catch (e: unknown) {
    const err = e as { response?: { data?: { error?: string } } };
    error.value = err?.response?.data?.error ?? 'Failed to join';
  }
}

async function handleStart() {
  await roomStore.startGame();
}

function handleLogout() {
  playerStore.logout();
  router.push('/auth');
}
</script>

<template>
  <!-- Header Bar -->
  <header class="header-bar">
    <span class="header-brand font-mono">TACTICAL TERMINAL</span>
    <div class="flex items-center gap-4">
      <button class="header-avatar" @click="handleLogout">
        <i class="pi pi-user"></i>
      </button>
    </div>
  </header>

  <div class="layout-terminal layout-centered">
    <!-- Hero -->
    <div class="hero-row">
      <div>
        <h1 class="hero-title font-display">{{ name || 'Game Lobby' }}</h1>
        <p class="hero-subtitle font-mono">
          Deployment Phase // Briefing Room
        </p>
      </div>
      <Button
        label="<< Back"
        class="btn-secondary-tactical"
        @click="router.push('/')"
      />
    </div>

    <!-- Error -->
    <div
      v-if="error"
      class="panel-low p-4 border-l-4 border-primary riveted mb-6"
    >
      <p class="text-xs font-mono text-primary mb-1">ERROR</p>
      <p class="text-sm font-mono">{{ error }}</p>
    </div>

    <!-- Deployment Status -->
    <section class="section-gap">
      <h2 class="section-heading font-mono">Deployment Status</h2>
      <LobbyStatus />
    </section>

    <!-- Role Selection -->
    <section class="section-gap">
      <template v-if="isGameMaster">
        <div v-if="gmPlayerRole === null">
          <h2 class="section-heading font-mono">Select Your Role</h2>
          <RoleSelector @select="handleRoleSelect" />
        </div>
        <div
          v-else
          class="panel-lowest p-4 border-l-4 border-secondary-container mb-6"
          style="max-width: 480px; margin: 0 auto"
        >
          <p class="text-xs font-mono text-secondary mb-1">
            GAME MASTER
          </p>
          <p class="text-sm font-display text-white">
            YOUR ROLE: {{ gmPlayerRole.toUpperCase() }}
          </p>
        </div>
      </template>

      <template v-else>
        <div v-if="role === null">
          <h2 class="section-heading font-mono">Select Your Role</h2>
          <RoleSelector @select="handleRoleSelect" />
        </div>
        <div
          v-else
          class="panel-low p-6 border-l-4 border-secondary-container riveted"
          style="max-width: 480px; margin: 0 auto"
        >
          <p class="text-xs font-mono text-secondary mb-1">
            ROLE ASSIGNED
          </p>
          <h2 class="text-xl font-display text-white mb-2">
            YOUR ROLE: {{ role.toUpperCase() }}
          </h2>
          <p class="text-sm text-surface-variant uppercase">
            Waiting for the Game Master to start the game...
          </p>
        </div>
      </template>
    </section>

    <!-- Army Selection (Placeholder) -->
    <section class="section-gap">
      <h2 class="section-heading font-mono">Army Selection</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Attacker Army -->
        <div class="army-slot">
          <div class="army-slot-header">
            <span class="army-slot-label font-mono">Attacker Army</span>
            <span class="army-slot-badge army-slot-badge--attacker font-mono">ATK</span>
          </div>
          <div class="army-slot-body">
            <div class="army-placeholder">
              <i class="pi pi-shield army-placeholder-icon"></i>
              <p class="font-mono text-sm">No army selected</p>
              <p class="font-mono army-hint">
                Choose a faction to deploy
              </p>
            </div>
          </div>
        </div>

        <!-- Defender Army -->
        <div class="army-slot">
          <div class="army-slot-header">
            <span class="army-slot-label font-mono">Defender Army</span>
            <span class="army-slot-badge army-slot-badge--defender font-mono">DEF</span>
          </div>
          <div class="army-slot-body">
            <div class="army-placeholder">
              <i class="pi pi-shield army-placeholder-icon"></i>
              <p class="font-mono text-sm">No army selected</p>
              <p class="font-mono army-hint">
                Choose a faction to deploy
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Start Game -->
    <section v-if="isGameMaster" style="max-width: 480px; margin: 0 auto">
      <Button
        label="Start Game"
        class="btn-tactical w-full h-16 text-xl"
        :disabled="!canStart"
        @click="handleStart"
      />
      <p
        v-if="!canStart"
        class="text-[10px] font-mono text-center mt-2 text-surface-variant uppercase"
      >
        Awaiting players to join before starting...
      </p>
    </section>
  </div>
</template>
