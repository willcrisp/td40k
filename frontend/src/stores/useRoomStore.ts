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

  // Whether the GM has also claimed a player slot
  const gmPlayerRole = computed<'attacker' | 'defender' | null>(() => {
    const player = usePlayerStore();
    const pid = player.playerId;
    if (!pid || gameMasterId.value !== pid) return null;
    if (attackerId.value === pid) return 'attacker';
    if (defenderId.value === pid) return 'defender';
    return null;
  });

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
    gmPlayerRole,
    canStart,
    loadRoom,
    applyServerState,
    joinRoom,
    startGame,
    nextPhase,
    prevPhase,
  };
});
