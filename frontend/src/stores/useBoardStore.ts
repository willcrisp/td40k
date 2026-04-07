import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useBoardStore = defineStore('board', () => {
  const zoom = ref(1.0);
  const panX = ref(0);
  const panY = ref(0);
  const rotation = ref(0); // 0, 90, 180, 270 degrees

  function reset() {
    zoom.value = 1.0;
    panX.value = 0;
    panY.value = 0;
    rotation.value = 0;
  }

  function toggleRotation() {
    rotation.value = (rotation.value + 90) % 360;
  }

  return { zoom, panX, panY, rotation, reset, toggleRotation };
});
