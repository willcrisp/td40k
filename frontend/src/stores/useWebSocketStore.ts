/// <reference types="vite/client" />
import { ref } from 'vue'
import { defineStore } from 'pinia'

type WsStatus = 'connecting' | 'connected' | 'disconnected'

/**
 * useWebSocketStore is the single, app-wide WebSocket connection manager.
 *
 * As a Pinia store (singleton), this instance is created ONCE and shared across
 * every component and store that calls useWebSocketStore(). This means there is
 * always exactly ONE WebSocket connection for the entire application, no matter
 * how many components are mounted.
 *
 * Other stores subscribe to incoming messages via `subscribe()`, receiving a
 * callback that fires each time the backend sends a JSON payload.
 */
export const useWebSocketStore = defineStore('websocket', () => {
  // ---------------------------------------------------------------------------
  // STATE
  // ---------------------------------------------------------------------------

  /** Reactive connection status — drives the live indicator in the UI */
  const status = ref<WsStatus>('connecting')

  /** The underlying native WebSocket instance — kept private to this store */
  let ws: WebSocket | null = null

  /**
   * Registry of subscriber callbacks. Any store/component can register a
   * function here; it will be invoked with every parsed JSON message from
   * the backend.
   */
  const subscribers = new Set<(data: unknown) => void>()

  // ---------------------------------------------------------------------------
  // ACTIONS
  // ---------------------------------------------------------------------------

  /**
   * Registers a callback that fires whenever a WebSocket message arrives.
   * Returns an unsubscribe function so callers can clean up if needed.
   */
  function subscribe(callback: (data: unknown) => void): () => void {
    subscribers.add(callback)
    return () => subscribers.delete(callback)
  }

  /**
   * Opens the WebSocket connection and wires up all event handlers.
   * Called once automatically at store creation; re-called on reconnect.
   */
  function connect() {
    // Dynamically pick ws:// or wss:// to match the page protocol
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'

    // In dev mode we bypass the Vite proxy and hit the backend port directly,
    // because some environments drop the WebSocket 101 upgrade through the proxy.
    const host = import.meta.env.DEV
      ? `${window.location.hostname}:8080`
      : window.location.host

    ws = new WebSocket(`${protocol}//${host}/api/ws`)

    ws.onopen = () => {
      status.value = 'connected'
    }

    ws.onmessage = (event: MessageEvent) => {
      try {
        const data: unknown = JSON.parse(event.data as string)
        // Fan out the message to every registered subscriber
        subscribers.forEach((cb) => cb(data))
      } catch (err) {
        console.error('[WebSocketStore] Failed to parse message:', err)
      }
    }

    ws.onerror = () => {
      status.value = 'disconnected'
    }

    ws.onclose = () => {
      status.value = 'disconnected'
      // Auto-reconnect after 3 s to survive backend restarts / network blips
      setTimeout(connect, 3000)
    }
  }

  /** Cleanly closes the socket (useful in tests or forced teardown) */
  function disconnect() {
    ws?.close()
    ws = null
  }

  // Open the connection immediately when the store is first instantiated
  connect()

  // ---------------------------------------------------------------------------
  // PUBLIC INTERFACE
  // ---------------------------------------------------------------------------
  return { status, subscribe, disconnect }
})
