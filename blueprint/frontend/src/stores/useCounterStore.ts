import { defineStore } from "pinia";
import { ref } from "vue";
import { apiGetCounter, apiIncrementCounter } from "@/lib/api";

export const useCounterStore = defineStore("counter", () => {
  const value = ref<number>(0);
  const loading = ref(false);

  async function fetchCounter() {
    const res = await apiGetCounter();
    value.value = res.data.value;
  }

  async function increment() {
    loading.value = true;
    try {
      await apiIncrementCounter();
      // value is updated via WebSocket broadcast
    } finally {
      loading.value = false;
    }
  }

  function applyUpdate(newValue: number) {
    value.value = newValue;
  }

  return { value, loading, fetchCounter, increment, applyUpdate };
});
