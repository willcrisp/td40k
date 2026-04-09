import { defineStore } from "pinia";
import { ref } from "vue";
import type { WsMessage, CounterUpdatePayload } from "@/types";
import { useCounterStore } from "./useCounterStore";

const WS_BASE = import.meta.env.VITE_WS_BASE_URL ?? "";

export const useWebSocketStore = defineStore("websocket", () => {
  const socket = ref<WebSocket | null>(null);
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  function connect() {
    if (socket.value?.readyState === WebSocket.OPEN) return;

    const url = `${WS_BASE}/ws`;
    const ws = new WebSocket(url);

    ws.onopen = () => {
      console.log("WebSocket connected");
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
      }
    };

    ws.onmessage = (event) => {
      try {
        const msg: WsMessage = JSON.parse(event.data as string);
        if (msg.event === "counter_update") {
          const payload = msg.payload as CounterUpdatePayload;
          useCounterStore().applyUpdate(payload.value);
        }
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected — reconnecting in 3s");
      reconnectTimer = setTimeout(connect, 3000);
    };

    ws.onerror = () => ws.close();

    socket.value = ws;
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
    socket.value?.close();
    socket.value = null;
  }

  return { socket, connect, disconnect };
});
