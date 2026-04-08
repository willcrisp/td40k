import { defineStore } from 'pinia';
import { ref } from 'vue';
import { v4 as uuidv4 } from 'uuid';

export const usePlayerStore = defineStore('player', () => {
  const playerId = ref<string>('');

  function init() {
    const stored = localStorage.getItem('player_id');
    if (stored) {
      playerId.value = stored;
    } else {
      const newId = uuidv4();
      playerId.value = newId;
      localStorage.setItem('player_id', newId);
    }
  }

  return {
    playerId,
    init,
  };
});
