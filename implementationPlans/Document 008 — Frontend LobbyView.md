Document 008 — Frontend: LobbyView

Purpose


Implement the waiting room where players select roles and the Game Master starts the game.


---

frontend/src/components/lobby/LobbyStatus.vue

	<script setup lang="ts">
	import { storeToRefs } from 'pinia';
	import Tag from 'primevue/tag';
	import { useRoomStore } from '@/stores/useRoomStore';
	
	const room = useRoomStore();
	const { attackerId, defenderId, name, status } = storeToRefs(room);
	</script>
	
	<template>
	  <div class="flex flex-col gap-3">
	    <h2 class="text-2xl font-bold">{{ name }}</h2>
	    <div class="flex gap-4">
	      <div class="flex items-center gap-2">
	        <span class="text-sm text-surface-400">Attacker:</span>
	        <Tag
	          :value="attackerId ? 'Ready' : 'Waiting…'"
	          :severity="attackerId ? 'success' : 'warn'"
	        />
	      </div>
	      <div class="flex items-center gap-2">
	        <span class="text-sm text-surface-400">Defender:</span>
	        <Tag
	          :value="defenderId ? 'Ready' : 'Waiting…'"
	          :severity="defenderId ? 'success' : 'warn'"
	        />
	      </div>
	    </div>
	    <p v-if="status === 'lobby'" class="text-sm text-surface-400">
	      Share this page URL for others to join.
	    </p>
	  </div>
	</template>


---

frontend/src/components/lobby/RoleSelector.vue

	<script setup lang="ts">
	import { storeToRefs } from 'pinia';
	import Button from 'primevue/button';
	import Card from 'primevue/card';
	import { useRoomStore } from '@/stores/useRoomStore';
	import { usePlayerStore } from '@/stores/usePlayerStore';
	
	const emit = defineEmits<{
	  (e: 'select', role: 'attacker' | 'defender'): void;
	}>();
	
	const room = useRoomStore();
	const player = usePlayerStore();
	const { attackerId, defenderId } = storeToRefs(room);
	const { playerId } = storeToRefs(player);
	
	function isMe(id: string | null) {
	  return id === playerId.value;
	}
	</script>
	
	<template>
	  <div class="grid grid-cols-2 gap-6 mt-6">
	    <!-- Attacker -->
	    <Card
	      :class="[
	        'border-2 transition-colors',
	        isMe(attackerId) ? 'border-red-500' : 'border-transparent',
	      ]"
	    >
	      <template #title>
	        <span class="text-red-400">⚔ Attacker</span>
	      </template>
	      <template #content>
	        <p class="text-sm text-surface-400 mb-4">
	          Controls the attacking force.
	        </p>
	        <div v-if="attackerId">
	          <span class="text-green-400 text-sm">
	            {{ isMe(attackerId) ? '✅ You' : '✅ Taken' }}
	          </span>
	        </div>
	        <Button
	          v-else
	          label="Choose Attacker"
	          severity="danger"
	          @click="emit('select', 'attacker')"
	        />
	      </template>
	    </Card>
	
	    <!-- Defender -->
	    <Card
	      :class="[
	        'border-2 transition-colors',
	        isMe(defenderId) ? 'border-blue-500' : 'border-transparent',
	      ]"
	    >
	      <template #title>
	        <span class="text-blue-400">🛡 Defender</span>
	      </template>
	      <template #content>
	        <p class="text-sm text-surface-400 mb-4">
	          Controls the defending force.
	        </p>
	        <div v-if="defenderId">
	          <span class="text-green-400 text-sm">
	            {{ isMe(defenderId) ? '✅ You' : '✅ Taken' }}
	          </span>
	        </div>
	        <Button
	          v-else
	          label="Choose Defender"
	          severity="info"
	          @click="emit('select', 'defender')"
	        />
	      </template>
	    </Card>
	  </div>
	</template>


---

frontend/src/views/LobbyView.vue

	<script setup lang="ts">
	import { onMounted, onUnmounted, watch } from 'vue';
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
	const { status, role, canStart, isGameMaster } = storeToRefs(roomStore);
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
	      <Message severity="info">
	        You are the Game Master. Wait for both players to join.
	      </Message>
	      <Button
	        label="Start Game"
	        class="mt-6 w-full"
	        size="large"
	        :disabled="!canStart"
	        @click="handleStart"
	      />
	    </div>
	
	    <!-- Player role selection -->
	    <div v-else-if="role === null">
	      <RoleSelector @select="handleRoleSelect" />
	    </div>
	
	    <!-- Already has a role -->
	    <div v-else class="mt-8">
	      <Message severity="success">
	        You've joined as
	        <strong>{{ role }}</strong>. Waiting for the Game Master to
	        start the game…
	      </Message>
	    </div>
	  </div>
	</template>


Note: Add import { ref } from 'vue'; to the script — it was omitted above for brevity but is required.