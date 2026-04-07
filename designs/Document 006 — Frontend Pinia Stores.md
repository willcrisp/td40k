Document 006 — Frontend: Pinia Stores

Purpose


Implement all Pinia stores with full TypeScript typing.


---

frontend/src/stores/usePlayerStore.ts

	import { defineStore } from 'pinia';
	import { ref } from 'vue';
	import { v4 as uuidv4 } from 'uuid';
	import { apiUpsertPlayer } from '@/lib/api';
	
	// Install uuid: bun add uuid && bun add -d @types/uuid
	
	export const usePlayerStore = defineStore('player', () => {
	  const playerId = ref<string>('');
	  const nickname = ref<string>('');
	  const initialized = ref(false);
	
	  async function initPlayer() {
	    let id = localStorage.getItem('player_id');
	    if (!id) {
	      id = uuidv4();
	      localStorage.setItem('player_id', id);
	    }
	    playerId.value = id;
	
	    const storedNick = localStorage.getItem('nickname') || '';
	    nickname.value = storedNick;
	    initialized.value = true;
	
	    // Only upsert if we have a nickname — HomeView prompts if not
	    if (storedNick) {
	      await apiUpsertPlayer(id, storedNick).catch(() => {});
	    }
	  }
	
	  async function setNickname(nick: string) {
	    nickname.value = nick;
	    localStorage.setItem('nickname', nick);
	    await apiUpsertPlayer(playerId.value, nick);
	  }
	
	  return { playerId, nickname, initialized, initPlayer, setNickname };
	});


---

frontend/src/stores/useGameListStore.ts

	import { defineStore } from 'pinia';
	import { ref } from 'vue';
	import type { OwnedGameSummary, JoinedGameSummary } from '@/types';
	import {
	  apiGetPlayerGames,
	  apiCreateRoom,
	  apiCloseRoom,
	} from '@/lib/api';
	import { usePlayerStore } from './usePlayerStore';
	
	export const useGameListStore = defineStore('gameList', () => {
	  const ownedGames = ref<OwnedGameSummary[]>([]);
	  const joinedGames = ref<JoinedGameSummary[]>([]);
	  const loading = ref(false);
	  const error = ref<string | null>(null);
	
	  async function fetchGames() {
	    const player = usePlayerStore();
	    loading.value = true;
	    error.value = null;
	    try {
	      const { data } = await apiGetPlayerGames(player.playerId);
	      ownedGames.value = data.owned;
	      joinedGames.value = data.joined;
	    } catch (e) {
	      error.value = 'Failed to load games';
	    } finally {
	      loading.value = false;
	    }
	  }
	
	  async function createGame(name: string): Promise<string> {
	    const { data } = await apiCreateRoom(name);
	    await fetchGames();
	    return data.id;
	  }
	
	  async function closeGame(roomId: string) {
	    await apiCloseRoom(roomId);
	    ownedGames.value = ownedGames.value.filter((g) => g.id !== roomId);
	  }
	
	  return {
	    ownedGames,
	    joinedGames,
	    loading,
	    error,
	    fetchGames,
	    createGame,
	    closeGame,
	  };
	});


---

frontend/src/stores/useRoomStore.ts

	import { defineStore } from 'pinia';
	import { ref, computed } from 'vue';
	import type {
	  Phase,
	  RoomStatus,
	  PlayerRole,
	  ActivePlayer,
	  RoomStatePayload,
	} from '@/types';
	import {
	  apiGetRoom,
	  apiJoinRoom,
	  apiStartGame,
	  apiPhaseNext,
	  apiPhasePrev,
	} from '@/lib/api';
	import { usePlayerStore } from './usePlayerStore';
	
	export const useRoomStore = defineStore('room', () => {
	  const roomId = ref<string | null>(null);
	  const name = ref<string>('');
	  const status = ref<RoomStatus>('lobby');
	  const gameMasterId = ref<string | null>(null);
	  const attackerId = ref<string | null>(null);
	  const defenderId = ref<string | null>(null);
	  const battleRound = ref<number>(1);
	  const activePlayer = ref<ActivePlayer>('attacker');
	  const currentPhase = ref<Phase>('command');
	  const winner = ref<ActivePlayer | null>(null);
	
	  const role = computed<PlayerRole>(() => {
	    const player = usePlayerStore();
	    const pid = player.playerId;
	    if (!pid || !roomId.value) return null;
	    if (gameMasterId.value === pid) return 'game_master';
	    if (attackerId.value === pid) return 'attacker';
	    if (defenderId.value === pid) return 'defender';
	    return null;
	  });
	
	  const isGameMaster = computed(() => role.value === 'game_master');
	
	  const canStart = computed(
	    () =>
	      status.value === 'lobby' &&
	      attackerId.value !== null &&
	      defenderId.value !== null &&
	      isGameMaster.value
	  );
	
	  async function loadRoom(id: string) {
	    const { data } = await apiGetRoom(id);
	    applyServerState({
	      room_id: data.id,
	      name: data.name,
	      status: data.status,
	      battle_round: data.battle_round,
	      active_player: data.active_player,
	      current_phase: data.current_phase,
	      winner: data.winner,
	      attacker_id: data.attacker_id,
	      defender_id: data.defender_id,
	      game_master_id: data.game_master_id,
	    });
	  }
	
	  function applyServerState(payload: RoomStatePayload) {
	    roomId.value = payload.room_id;
	    name.value = payload.name;
	    status.value = payload.status;
	    battleRound.value = payload.battle_round;
	    activePlayer.value = payload.active_player;
	    currentPhase.value = payload.current_phase;
	    winner.value = payload.winner;
	    attackerId.value = payload.attacker_id;
	    defenderId.value = payload.defender_id;
	    gameMasterId.value = payload.game_master_id;
	  }
	
	  async function joinRoom(
	    id: string,
	    selectedRole: 'attacker' | 'defender'
	  ) {
	    await apiJoinRoom(id, selectedRole);
	  }
	
	  async function startGame() {
	    if (!roomId.value) return;
	    await apiStartGame(roomId.value);
	  }
	
	  async function nextPhase() {
	    if (!roomId.value) return;
	    await apiPhaseNext(roomId.value);
	  }
	
	  async function prevPhase() {
	    if (!roomId.value) return;
	    await apiPhasePrev(roomId.value);
	  }
	
	  return {
	    roomId,
	    name,
	    status,
	    gameMasterId,
	    attackerId,
	    defenderId,
	    battleRound,
	    activePlayer,
	    currentPhase,
	    winner,
	    role,
	    isGameMaster,
	    canStart,
	    loadRoom,
	    applyServerState,
	    joinRoom,
	    startGame,
	    nextPhase,
	    prevPhase,
	  };
	});


---

frontend/src/stores/useBoardStore.ts

	import { defineStore } from 'pinia';
	import { ref } from 'vue';
	
	export const useBoardStore = defineStore('board', () => {
	  const zoom = ref(1.0);
	  const panX = ref(0);
	  const panY = ref(0);
	
	  function reset() {
	    zoom.value = 1.0;
	    panX.value = 0;
	    panY.value = 0;
	  }
	
	  return { zoom, panX, panY, reset };
	});


---

frontend/src/stores/useWebSocketStore.ts

	import { defineStore } from 'pinia';
	import { ref } from 'vue';
	import type { WsMessage } from '@/types';
	import { useRoomStore } from './useRoomStore';
	
	const WS_BASE =
	  import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080';
	
	const RECONNECT_DELAY_MS = 3000;
	
	export const useWebSocketStore = defineStore('websocket', () => {
	  const connected = ref(false);
	  let socket: WebSocket | null = null;
	  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	  let currentRoomId = '';
	  let currentPlayerId = '';
	
	  function connect(roomId: string, playerId: string) {
	    currentRoomId = roomId;
	    currentPlayerId = playerId;
	    _open();
	  }
	
	  function disconnect() {
	    if (reconnectTimer) clearTimeout(reconnectTimer);
	    if (socket) {
	      socket.onclose = null; // suppress reconnect on manual close
	      socket.close();
	      socket = null;
	    }
	    connected.value = false;
	  }
	
	  function _open() {
	    const url = `${WS_BASE}/ws?room_id=${currentRoomId}&player_id=${currentPlayerId}`;
	    socket = new WebSocket(url);
	
	    socket.onopen = () => {
	      connected.value = true;
	    };
	
	    socket.onmessage = (event: MessageEvent) => {
	      try {
	        const msg: WsMessage = JSON.parse(event.data);
	        if (msg.event === 'room_state') {
	          const roomStore = useRoomStore();
	          roomStore.applyServerState(msg.payload);
	        }
	      } catch {
	        // ignore malformed messages
	      }
	    };
	
	    socket.onclose = () => {
	      connected.value = false;
	      reconnectTimer = setTimeout(_open, RECONNECT_DELAY_MS);
	    };
	
	    socket.onerror = () => {
	      socket?.close();
	    };
	  }
	
	  return { connected, connect, disconnect };
	});