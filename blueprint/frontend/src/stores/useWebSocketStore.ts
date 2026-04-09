import { defineStore } from "pinia";
import { ref } from "vue";
import type { WsMessage, CounterUpdatePayload, NoteEvent } from "@/types";
import { useCounterStore } from "./useCounterStore";
import { useNotesStore } from "./useNotesStore";

const WS_BASE = import.meta.env.VITE_WS_BASE_URL ?? "";

const MIN_BACKOFF = 1000;
const MAX_BACKOFF = 30000;

export const useWebSocketStore = defineStore("websocket", () => {
  const socket = ref<WebSocket | null>(null);
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let backoff = MIN_BACKOFF;

  function connect() {
    if (socket.value?.readyState === WebSocket.OPEN) return;

    const ws = new WebSocket(`${WS_BASE}/ws`);

    ws.onopen = () => {
      backoff = MIN_BACKOFF;
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
      }
    };

    ws.onmessage = (event) => {
      try {
        const msg: WsMessage = JSON.parse(event.data as string);

        if (msg.event === "counter_update") {
          const p = msg.payload as CounterUpdatePayload;
          useCounterStore().applyUpdate(p.value);
        } else if (msg.event === "notes_update") {
          const p = msg.payload as NoteEvent;
          const notesStore = useNotesStore();
          if (p.op === "insert" && p.player_id && p.content && p.created_at) {
            notesStore.applyInsert({
              id: p.id,
              player_id: p.player_id,
              username: p.username ?? "",
              content: p.content,
              created_at: p.created_at,
            });
          } else if (p.op === "delete") {
            notesStore.applyDelete(p.id);
          }
        }
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = () => {
      reconnectTimer = setTimeout(() => {
        connect();
      }, backoff);
      backoff = Math.min(backoff * 2, MAX_BACKOFF);
    };

    ws.onerror = () => ws.close();

    socket.value = ws;
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
    backoff = MIN_BACKOFF;
    socket.value?.close();
    socket.value = null;
  }

  return { socket, connect, disconnect };
});
