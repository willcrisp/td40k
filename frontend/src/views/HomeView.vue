<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
import ConfirmDialog from 'primevue/confirmdialog';
import { useConfirm } from 'primevue/useconfirm';
import { useGameListStore } from '@/stores/useGameListStore';
import { usePlayerStore } from '@/stores/usePlayerStore';
import CreateGameModal from '@/components/home/CreateGameModal.vue';
import OwnedGameCard from '@/components/home/OwnedGameCard.vue';
import JoinedGameCard from '@/components/home/JoinedGameCard.vue';

const router = useRouter();
const gameListStore = useGameListStore();
const playerStore = usePlayerStore();
const confirm = useConfirm();

const { ownedGames, joinedGames, loading } = storeToRefs(gameListStore);
const { nickname, username } = storeToRefs(playerStore);

import { ref } from 'vue';
const showCreateModal = ref(false);

onMounted(async () => {
  await gameListStore.fetchGames();
});

async function handleCreate(name: string) {
  showCreateModal.value = false;
  const id = await gameListStore.createGame(name);
  router.push(`/lobby/${id}`);
}

function handleEnterGame(id: string) {
  router.push(`/lobby/${id}`);
}

function handleCloseGame(id: string) {
  confirm.require({
    message:
      'Are you sure you want to close this game? ' +
      'All players will be disconnected.',
    header: 'Close Game',
    icon: 'pi pi-exclamation-triangle',
    acceptSeverity: 'danger',
    accept: () => gameListStore.closeGame(id),
  });
}

function handleRejoin(id: string) {
  router.push(`/lobby/${id}`);
}

function handleLogout() {
  playerStore.logout();
  router.push('/auth');
}
</script>

<template>
  <ConfirmDialog />

  <div class="layout-terminal">
    <!-- Header -->
    <div class="flex items-start justify-between mb-8">
      <div>
        <h1 class="text-3xl font-display text-primary">Game Tracker</h1>
        <p class="text-sm font-mono text-surface-variant">
          Status: Online // Sector 40K
        </p>
      </div>
      <div class="flex items-center gap-4">
        <div class="text-right">
          <p class="text-xs font-mono text-tertiary">Logged In As</p>
          <p class="text-sm font-display">{{ nickname || username }}</p>
        </div>
        <Button
          label="Logout"
          severity="secondary"
          size="small"
          class="font-mono"
          @click="handleLogout"
        />
      </div>
    </div>

    <!-- Owned games -->
    <section class="section-gap">
      <div
        class="flex items-end justify-between mb-6 border-b border-ghost-border pb-2"
      >
        <div>
          <h2 class="text-2xl font-display">My Games</h2>
          <p class="text-sm font-mono text-surface-variant">
            Games you are managing
          </p>
        </div>
        <Button
          label="+ New Game"
          @click="showCreateModal = true"
          class="btn-tactical"
        />
      </div>

      <div v-if="loading" class="font-mono text-tertiary">
        Loading games...
      </div>

      <div
        v-else-if="ownedGames.length"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        <OwnedGameCard
          v-for="game in ownedGames"
          :key="game.id"
          :game="game"
          @enter="handleEnterGame"
          @close="handleCloseGame"
        />
      </div>
      <div
        v-else
        class="panel-lowest p-6 border-l-4 border-outline-variant"
      >
        <p class="text-surface-variant font-mono">
          No active games found. Create one to get started.
        </p>
      </div>
    </section>

    <!-- Joined games -->
    <section>
      <div class="mb-6 border-b border-ghost-border pb-2">
        <h2 class="text-2xl font-display">Joined Games</h2>
        <p class="text-sm font-mono text-surface-variant">
          Games where you are a player
        </p>
      </div>

      <div
        v-if="joinedGames.length"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        <JoinedGameCard
          v-for="game in joinedGames"
          :key="game.id"
          :game="game"
          @rejoin="handleRejoin"
        />
      </div>
      <div
        v-else
        class="panel-lowest p-6 border-l-4 border-outline-variant"
      >
        <p class="text-surface-variant font-mono">
          You haven't joined any games yet.
        </p>
      </div>
    </section>
  </div>

  <CreateGameModal
    v-model:visible="showCreateModal"
    @create="handleCreate"
  />
</template>
