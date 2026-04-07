<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
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
const { nickname } = storeToRefs(playerStore);

const showCreateModal = ref(false);
const nickInput = ref('');
const needsNickname = computed(() => !nickname.value);

onMounted(async () => {
  if (nickname.value) {
    await gameListStore.fetchGames();
  }
});

async function saveNickname() {
  if (!nickInput.value.trim()) return;
  await playerStore.setNickname(nickInput.value.trim());
  await gameListStore.fetchGames();
}

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
</script>

<template>
  <ConfirmDialog />

  <div class="min-h-screen p-6 max-w-5xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-3xl font-bold tracking-widest uppercase">
        T40K
      </h1>
      <span v-if="nickname" class="text-sm text-surface-400">
        Commander: {{ nickname }}
      </span>
    </div>

    <!-- Nickname prompt -->
    <div
      v-if="needsNickname"
      class="flex flex-col items-center gap-4 py-20"
    >
      <h2 class="text-xl font-semibold">
        Enter your commander name to begin
      </h2>
      <div class="flex gap-2 w-80">
        <InputText
          v-model="nickInput"
          placeholder="e.g. Valdris"
          class="flex-1"
          @keyup.enter="saveNickname"
          autofocus
        />
        <Button label="Enter" @click="saveNickname" />
      </div>
    </div>

    <!-- Game lists -->
    <template v-else>
      <!-- Owned games -->
      <section class="mb-10">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-xl font-semibold uppercase tracking-wide">
            My Games
            <span class="text-surface-400 text-sm ml-2">(Game Master)</span>
          </h2>
          <Button
            label="+ Create Game"
            @click="showCreateModal = true"
          />
        </div>

        <div v-if="loading" class="text-surface-400">Loading…</div>

        <div
          v-else-if="ownedGames.length"
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
        >
          <OwnedGameCard
            v-for="game in ownedGames"
            :key="game.id"
            :game="game"
            @enter="handleEnterGame"
            @close="handleCloseGame"
          />
        </div>
        <p v-else class="text-surface-400">
          No games yet. Create one to get started.
        </p>
      </section>

      <!-- Joined games -->
      <section>
        <h2 class="text-xl font-semibold uppercase tracking-wide mb-4">
          Games I've Joined
        </h2>

        <div
          v-if="joinedGames.length"
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
        >
          <JoinedGameCard
            v-for="game in joinedGames"
            :key="game.id"
            :game="game"
            @rejoin="handleRejoin"
          />
        </div>
        <p v-else class="text-surface-400">
          You haven't joined any games yet.
        </p>
      </section>
    </template>
  </div>

  <CreateGameModal
    v-model:visible="showCreateModal"
    @create="handleCreate"
  />
</template>
