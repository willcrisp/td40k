Document 005 — Frontend: Project Setup & Global Config

Purpose


Bootstrap the Vue 3 frontend project with all dependencies, global config, types, the API client, and the router.


---

Initialize Project

	cd frontend
	bun create vite . --template vue-ts
	bun install
	bun add pinia vue-router @primevue/themes primevue primeicons axios
	bun add -d prettier @vue/tsconfig


---

frontend/package.json (scripts section)

	{
	  "scripts": {
	    "dev": "vite",
	    "build": "vue-tsc && vite build",
	    "preview": "vite preview",
	    "format": "prettier --write src/"
	  }
	}


---

frontend/.prettierrc

	{
	  "printWidth": 80,
	  "singleQuote": true,
	  "semi": true,
	  "trailingComma": "es5"
	}


---

frontend/vite.config.ts

	import { defineConfig } from 'vite';
	import vue from '@vitejs/plugin-vue';
	import { resolve } from 'path';
	
	export default defineConfig({
	  plugins: [vue()],
	  resolve: {
	    alias: {
	      '@': resolve(__dirname, 'src'),
	    },
	  },
	  server: {
	    port: 5173,
	    proxy: {
	      '/api': {
	        target: process.env.VITE_API_BASE_URL || 'http://localhost:8080',
	        changeOrigin: true,
	      },
	      '/ws': {
	        target: process.env.VITE_WS_BASE_URL || 'ws://localhost:8080',
	        ws: true,
	        changeOrigin: true,
	      },
	    },
	  },
	});


---

frontend/src/types/index.ts

	export type Phase =
	  | 'command'
	  | 'movement'
	  | 'shooting'
	  | 'charge'
	  | 'fight';
	
	export type RoomStatus = 'lobby' | 'active' | 'finished' | 'closed';
	
	export type PlayerRole =
	  | 'attacker'
	  | 'defender'
	  | 'game_master'
	  | null;
	
	export type ActivePlayer = 'attacker' | 'defender';
	
	export interface Room {
	  id: string;
	  name: string;
	  status: RoomStatus;
	  game_master_id: string;
	  attacker_id: string | null;
	  defender_id: string | null;
	  battle_round: number;
	  active_player: ActivePlayer;
	  current_phase: Phase;
	  winner: ActivePlayer | null;
	  created_at: string;
	  updated_at: string;
	}
	
	export interface OwnedGameSummary {
	  id: string;
	  name: string;
	  status: RoomStatus;
	  battle_round: number;
	  active_player: ActivePlayer;
	  current_phase: Phase;
	  attacker_id: string | null;
	  defender_id: string | null;
	  created_at: string;
	}
	
	export interface JoinedGameSummary {
	  id: string;
	  name: string;
	  status: RoomStatus;
	  role: 'attacker' | 'defender';
	  battle_round: number;
	  current_phase: Phase;
	  created_at: string;
	}
	
	export interface RoomStatePayload {
	  room_id: string;
	  name: string;
	  status: RoomStatus;
	  battle_round: number;
	  active_player: ActivePlayer;
	  current_phase: Phase;
	  winner: ActivePlayer | null;
	  attacker_id: string | null;
	  defender_id: string | null;
	  game_master_id: string;
	}
	
	export interface WsMessage {
	  event: 'room_state';
	  payload: RoomStatePayload;
	}
	
	export const PHASES: Phase[] = [
	  'command',
	  'movement',
	  'shooting',
	  'charge',
	  'fight',
	];
	
	export const PHASE_LABELS: Record<Phase, string> = {
	  command: 'Command Phase',
	  movement: 'Movement Phase',
	  shooting: 'Shooting Phase',
	  charge: 'Charge Phase',
	  fight: 'Fight Phase',
	};
	
	export const PHASE_NUMBERS: Record<Phase, number> = {
	  command: 1,
	  movement: 2,
	  shooting: 3,
	  charge: 4,
	  fight: 5,
	};


---

frontend/src/lib/api.ts

	import axios from 'axios';
	
	const BASE = import.meta.env.VITE_API_BASE_URL || '';
	
	const client = axios.create({ baseURL: BASE });
	
	// Inject X-Player-ID on every request
	client.interceptors.request.use((config) => {
	  const playerId = localStorage.getItem('player_id');
	  if (playerId) {
	    config.headers['X-Player-ID'] = playerId;
	  }
	  return config;
	});
	
	export const apiUpsertPlayer = (id: string, nickname: string) =>
	  client.post('/api/players', { id, nickname });
	
	export const apiGetPlayerGames = (id: string) =>
	  client.get(`/api/players/${id}/games`);
	
	export const apiCreateRoom = (name: string) =>
	  client.post<{ id: string }>('/api/rooms', { name });
	
	export const apiGetRoom = (id: string) =>
	  client.get(`/api/rooms/${id}`);
	
	export const apiJoinRoom = (
	  roomId: string,
	  role: 'attacker' | 'defender'
	) => client.post(`/api/rooms/${roomId}/join`, { role });
	
	export const apiStartGame = (roomId: string) =>
	  client.post(`/api/rooms/${roomId}/start`);
	
	export const apiPhaseNext = (roomId: string) =>
	  client.post(`/api/rooms/${roomId}/phase/next`);
	
	export const apiPhasePrev = (roomId: string) =>
	  client.post(`/api/rooms/${roomId}/phase/prev`);
	
	export const apiCloseRoom = (roomId: string) =>
	  client.post(`/api/rooms/${roomId}/close`);


---

frontend/src/router/index.ts

	import {
	  createRouter,
	  createWebHistory,
	  type NavigationGuardNext,
	  type RouteLocationNormalized,
	} from 'vue-router';
	import HomeView from '@/views/HomeView.vue';
	import LobbyView from '@/views/LobbyView.vue';
	import GameView from '@/views/GameView.vue';
	import { usePlayerStore } from '@/stores/usePlayerStore';
	import { useRoomStore } from '@/stores/useRoomStore';
	
	const router = createRouter({
	  history: createWebHistory(),
	  routes: [
	    {
	      path: '/',
	      component: HomeView,
	    },
	    {
	      path: '/lobby/:id',
	      component: LobbyView,
	      beforeEnter: requirePlayer,
	    },
	    {
	      path: '/game/:id',
	      component: GameView,
	      beforeEnter: [requirePlayer, requireParticipant],
	    },
	  ],
	});
	
	function requirePlayer(
	  _to: RouteLocationNormalized,
	  _from: RouteLocationNormalized,
	  next: NavigationGuardNext
	) {
	  const player = usePlayerStore();
	  if (!player.playerId) return next('/');
	  next();
	}
	
	async function requireParticipant(
	  to: RouteLocationNormalized,
	  _from: RouteLocationNormalized,
	  next: NavigationGuardNext
	) {
	  const player = usePlayerStore();
	  const room = useRoomStore();
	  await room.loadRoom(to.params.id as string);
	
	  const isParticipant =
	    room.gameMasterId === player.playerId ||
	    room.attackerId === player.playerId ||
	    room.defenderId === player.playerId;
	
	  if (!isParticipant) return next(`/lobby/${to.params.id}`);
	  if (room.status === 'lobby') return next(`/lobby/${to.params.id}`);
	  next();
	}
	
	export default router;


---

frontend/src/main.ts

	import { createApp } from 'vue';
	import { createPinia } from 'pinia';
	import PrimeVue from 'primevue/config';
	import Aura from '@primevue/themes/aura';
	import 'primeicons/primeicons.css';
	import router from '@/router';
	import App from './App.vue';
	
	const app = createApp(App);
	
	app.use(createPinia());
	app.use(router);
	app.use(PrimeVue, {
	  theme: {
	    preset: Aura,
	    options: {
	      darkModeSelector: '.dark',
	    },
	  },
	});
	
	app.mount('#app');


---

frontend/src/App.vue

	<script setup lang="ts">
	import { onMounted } from 'vue';
	import { usePlayerStore } from '@/stores/usePlayerStore';
	
	const playerStore = usePlayerStore();
	
	onMounted(async () => {
	  await playerStore.initPlayer();
	});
	</script>
	
	<template>
	  <div class="dark">
	    <RouterView />
	  </div>
	</template>


---