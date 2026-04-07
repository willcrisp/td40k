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
