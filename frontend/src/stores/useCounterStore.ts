import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useWebSocketStore } from './useWebSocketStore'

// ---------------------------------------------------------------------------
// TYPES
// ---------------------------------------------------------------------------

/**
 * The exact shape of the JSON payload the Go backend sends over WebSocket
 * (and also returns from the REST `/api/counter` endpoint).
 */
export interface CounterResponse {
  name?: string
  value: number
  updated_at: string
}

/**
 * useCounterStore owns all state and business logic for the counter domain.
 *
 * It is the Pinia equivalent of the old `useCounter` composable, but with
 * a crucial advantage: because it is a store, its state is a TRUE singleton.
 * Any component that calls `useCounterStore()` receives the same reactive
 * instance — no duplicated fetches, no duplicated WebSocket subscriptions.
 */
export const useCounterStore = defineStore('counter', () => {
  // ---------------------------------------------------------------------------
  // DEPENDENCIES
  // ---------------------------------------------------------------------------

  /**
   * We pull in the WebSocket store to read connection status and subscribe to
   * incoming messages. The WS store is itself a singleton, so this is safe
   * to call from any number of components.
   */
  const wsStore = useWebSocketStore()

  // ---------------------------------------------------------------------------
  // STATE
  // ---------------------------------------------------------------------------

  /** The name of the counter "room" the user is currently viewing */
  const counterName = ref('main')

  /** The live counter value received from the database — null while loading */
  const counter = ref<number | null>(null)

  /** ISO timestamp of the last database write, used in the meta row */
  const updatedAt = ref<string | null>(null)

  /** True while the initial REST fetch is in-flight */
  const loading = ref(true)

  /** True while the increment POST request is in-flight (prevents double-clicks) */
  const incrementing = ref(false)

  /** Non-null when a network or parsing error has occurred */
  const error = ref<string | null>(null)

  /**
   * Toggled to `true` for 600 ms whenever a WebSocket push arrives.
   * Drives the CSS `pop` keyframe animation on the counter value.
   */
  const wssAnimating = ref(false)

  /** Debounce timer handle — prevents firing a fetch on every keystroke */
  let fetchTimeout: ReturnType<typeof setTimeout> | null = null

  // ---------------------------------------------------------------------------
  // COMPUTED
  // ---------------------------------------------------------------------------

  /**
   * Proxy the WebSocket connection status into this store so that components
   * only need to import a single store.
   */
  const connectionStatus = computed(() => wsStore.status)

  // ---------------------------------------------------------------------------
  // WEBSOCKET SUBSCRIPTION
  // ---------------------------------------------------------------------------

  /**
   * Register our message handler with the WebSocket store.
   * The returned unsubscribe function is intentionally not called here because
   * this is a singleton store that lives for the entire app lifetime.
   */
  wsStore.subscribe((raw: unknown) => {
    const data = raw as CounterResponse

    // FILTRATION LAYER:
    // Postgres broadcasts ALL counter updates to ALL connected clients.
    // We must ignore updates that belong to a different counter room.
    const targetName = counterName.value.trim() || 'main'
    if (data.name && data.name !== targetName) return

    // Apply the new values — Vue's reactivity system re-renders any bound DOM instantly
    counter.value = data.value
    updatedAt.value = data.updated_at

    // Trigger the pop animation for 600 ms
    wssAnimating.value = true
    setTimeout(() => { wssAnimating.value = false }, 600)
  })

  // ---------------------------------------------------------------------------
  // ACTIONS
  // ---------------------------------------------------------------------------

  /**
   * Fetches the current counter value from the REST API.
   * Called on mount and whenever the user switches to a new counter name.
   */
  async function fetchCounter() {
    const name = counterName.value.trim() || 'main'
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`/api/counter?name=${encodeURIComponent(name)}`)
      if (!res.ok) throw new Error(`Failed to fetch counter: ${name}`)
      const data: CounterResponse = await res.json()
      counter.value = data.value
      updatedAt.value = data.updated_at
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  /**
   * Debounced handler for the counter name input.
   * Waits 400 ms of silence before firing a new fetch to avoid network spam.
   */
  function onNameInput() {
    if (fetchTimeout) clearTimeout(fetchTimeout)
    if (!counterName.value.trim()) return
    loading.value = true
    fetchTimeout = setTimeout(() => { fetchCounter() }, 400)
  }

  /**
   * Sends a POST to the backend to increment the named counter.
   * The real-time update arrives back via WebSocket, so we intentionally do
   * NOT update `counter` here — the WS push is the single source of truth.
   */
  async function increment() {
    if (incrementing.value) return // Guard against rapid double-clicks
    incrementing.value = true
    error.value = null
    try {
      const res = await fetch('/api/counter/increment', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: counterName.value.trim() || 'main' }),
      })
      if (!res.ok) throw new Error('Failed to mutate counter')
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      incrementing.value = false
    }
  }

  // Kick off the initial data load as soon as the store is instantiated
  fetchCounter()

  // ---------------------------------------------------------------------------
  // PUBLIC INTERFACE
  // ---------------------------------------------------------------------------
  return {
    // State
    counterName,
    counter,
    updatedAt,
    loading,
    incrementing,
    error,
    wssAnimating,
    // Computed
    connectionStatus,
    // Actions
    fetchCounter,
    onNameInput,
    increment,
  }
})
