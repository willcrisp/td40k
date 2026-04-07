Document 007 — Frontend: HomeView

Purpose


Implement the landing page with owned games, joined games, nickname prompt, and create game modal.


---

frontend/src/components/home/CreateGameModal.vue

	<script setup lang="ts">
	import { ref } from 'vue';
	import Dialog from 'primevue/dialog';
	import InputText from 'primevue/inputtext';
	import Button from 'primevue/button';
	
	const props = defineProps<{ visible: boolean }>();
	const emit = defineEmits<{
	  (e: 'update:visible', v: boolean): void;
	  (e: 'create', name: string): void;
	}>();
	
	const gameName = ref('');
	
	function submit() {
	  if (!gameName.value.trim()) return;
	  emit('create', gameName.value.trim());
	  gameName.value = '';
	}
	</script>
	
	<template>
	  <Dialog
	    :visible="props.visible"
	    @update:visible="emit('update:visible', $event)"
	    header="Create New Game"
	    modal
	    :style="{ width: '400px' }"
	  >
	    <div class="flex flex-col gap-4 pt-2">
	      <label for="game-name" class="font-semibold">Game Name</label>
	      <InputText
	        id="game-name"
	        v-model="gameName"
	        placeholder="e.g. Battle for Cadia"
	        @keyup.enter="submit"
	        autofocus
	      />
	    </div>
	    <template #footer>
	      <Button
	        label="Cancel"
	        severity="secondary"
	        @click="emit('update:visible', false)"
	      />
	      <Button
	        label="Create"
	        :disabled="!gameName.trim()"
	        @click="submit"
	      />
	    </template>
	  </Dialog>
	</template>


---

frontend/src/components/home/OwnedGameCard.vue

	<script setup lang="ts">
	import type { OwnedGameSummary } from '@/types';
	import { PHASE_LABELS } from '@/types';
	import Button from 'primevue/button';
	import Card from 'primevue/card';
	import Tag from 'primevue/tag';
	
	const props = defineProps<{ game: OwnedGameSummary }>();
	const emit = defineEmits<{
	  (e: 'enter', id: string): void;
	  (e: 'close', id: string): void;
	}>();
	
	const statusSeverity: Record<string, string> = {
	  lobby: 'warn',
	  active: 'success',
	  finished: 'info',
	  closed: 'secondary',
	};
	</script>
	
	<template>
	  <Card>
	    <template #title>{{ props.game.name }}</template>
	    <template #subtitle>
	      <Tag
	        :value="props.game.status.toUpperCase()"
	        :severity="statusSeverity[props.game.status]"
	      />
	    </template>
	    <template #content>
	      <p v-if="props.game.status === 'active'" class="text-sm">
	        Round {{ props.game.battle_round }} of 5 &mdash;
	        {{ props.game.active_player.toUpperCase() }}&apos;s Turn
	        <br />
	        {{ PHASE_LABELS[props.game.current_phase] }}
	      </p>
	      <p v-else-if="props.game.status === 'lobby'" class="text-sm">
	        Waiting for players&hellip;
	        <br />
	        Attacker:
	        {{ props.game.attacker_id ? '✅ Joined' : '⏳ Empty' }}
	        &nbsp;|&nbsp; Defender:
	        {{ props.game.defender_id ? '✅ Joined' : '⏳ Empty' }}
	      </p>
	      <p v-else-if="props.game.status === 'finished'" class="text-sm">
	        Game complete
	      </p>
	    </template>
	    <template #footer>
	      <div class="flex gap-2">
	        <Button
	          label="Enter"
	          size="small"
	          @click="emit('enter', props.game.id)"
	        />
	        <Button
	          label="Close Game"
	          size="small"
	          severity="danger"
	          outlined
	          :disabled="
	            props.game.status === 'closed' ||
	            props.game.status === 'finished'
	          "
	          @click="emit('close', props.game.id)"
	        />
	      </div>
	    </template>
	  </Card>
	</template>


---

frontend/src/components/home/JoinedGameCard.vue

	<script setup lang="ts">
	import type { JoinedGameSummary } from '@/types';
	import { PHASE_LABELS } from '@/types';
	import Button from 'primevue/button';
	import Card from 'primevue/card';
	import Tag from 'primevue/tag';
	
	const props = defineProps<{ game: JoinedGameSummary }>();
	const emit = defineEmits<{
	  (e: 'rejoin', id: string): void;
	}>();
	
	const statusSeverity: Record<string, string> = {
	  lobby: 'warn',
	  active: 'success',
	  finished: 'info',
	  closed: 'secondary',
	};
	
	const roleLabel: Record<string, string> = {
	  attacker: 'Attacker',
	  defender: 'Defender',
	};
	</script>
	
	<template>
	  <Card>
	    <template #title>{{ props.game.name }}</template>
	    <template #subtitle>
	      <div class="flex gap-2">
	        <Tag
	          :value="props.game.status.toUpperCase()"
	          :severity="statusSeverity[props.game.status]"
	        />
	        <Tag :value="roleLabel[props.game.role]" severity="info" />
	      </div>
	    </template>
	    <template #content>
	      <p class="text-sm">
	        Round {{ props.game.battle_round }} of 5
	        <br />
	        {{ PHASE_LABELS[props.game.current_phase] }}
	      </p>
	    </template>
	    <template #footer>
	      <Button
	        label="Rejoin"
	        size="small"
	        :disabled="props.game.status === 'closed'"
	        @click="emit('rejoin', props.game.id)"
	      />
	    </template>
	  </Card>
	</template>


---

frontend/src/views/HomeView.vue

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