import { defineStore } from 'pinia';
import { ref } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import { apiUpsertPlayer } from '@/lib/api';

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
