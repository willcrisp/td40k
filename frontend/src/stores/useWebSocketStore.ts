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
