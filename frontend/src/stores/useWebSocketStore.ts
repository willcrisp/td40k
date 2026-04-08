import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { WsMessage, GameUnitsUpdate } from '@/types';
import { useRoomStore } from './useRoomStore';
import { useUnitStore } from './useUnitStore';

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
        const msg = JSON.parse(event.data);
        if (msg.event === 'room_state') {
          const roomStore = useRoomStore();
          roomStore.applyServerState(msg.payload);
        } else if (msg.event === 'game_units_updates') {
          const unitStore = useUnitStore();
          const payload = msg.payload;

          if (payload.event_type === 'unit_removed') {
            unitStore.removeUnit(payload.unit_id);
          } else {
            // For unit_placed and unit_moved, we fetch the updated unit
            // This will be handled by the component that cares about units
            // For now, we just update local cache
            unitStore.updateUnit(payload.unit_id, {
              x: payload.x,
              y: payload.y,
              facing_degrees: payload.facing_degrees,
              status: payload.status,
              wounds: payload.wounds,
            });
          }
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
