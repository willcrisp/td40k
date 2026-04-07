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
