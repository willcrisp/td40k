<script setup lang="ts">
// PrimeVue UI framework UI imports
import Card from 'primevue/card'
import Button from 'primevue/button'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import InputText from 'primevue/inputtext'
import FloatLabel from 'primevue/floatlabel'
import { storeToRefs } from 'pinia'

// Import the Pinia store instead of the old composable.
// useCounterStore() always returns the SAME singleton instance, so multiple
// components will share state without any extra wiring.
import { useCounterStore } from '../../../stores/useCounterStore'

const counterStore = useCounterStore()

// storeToRefs extracts reactive refs from the store while keeping their
// reactivity intact — the Pinia equivalent of destructuring a composable.
const {
  counterName, counter, updatedAt, loading, incrementing,
  error, wssAnimating, connectionStatus,
} = storeToRefs(counterStore)

// Actions are plain functions — destructure them directly (no storeToRefs needed)
const { onNameInput, increment } = counterStore

// Local formatting utility
function formatDate(iso: string | null) {
  if (!iso) return '—'
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'medium',
  }).format(new Date(iso))
}
</script>

<template>
  <Card class="counter-card">
    <template #content>
      
      <!-- DYNAMIC ROOM SELECTOR -->
      <!-- 'v-model' creates a two-way data binding. As you type, 'counterName' in the store updates automatically! -->
      <!-- '@update:modelValue' fires our debouncer whenever the text shifts -->
      <div class="counter-selector">
        <FloatLabel>
          <InputText id="counter-name" v-model="counterName" @update:modelValue="onNameInput" autocomplete="off" />
          <label for="counter-name">Active Counter Name</label>
        </FloatLabel>
      </div>

      <!-- If waiting on initial load, show a PrimeVue spinner -->
      <div v-if="loading" class="spinner-wrap">
        <ProgressSpinner style="width: 48px; height: 48px" strokeWidth="4" />
      </div>

      <!-- MAIN LOADED UI BODY -->
      <div v-else class="counter-body">
        
        <!-- Live status indicator tracking the socket pulse! -->
        <div class="ws-status">
          <span class="dot" :class="connectionStatus"></span>
          <span class="status-text">{{ connectionStatus === 'connected' ? 'Live Sync Active' : 'Connecting...' }}</span>
        </div>

        <!-- The actual Counter metric! NOTE: ':class' evaluates an expression. When 'pop' is true, it triggers CSS keyframes -->
        <div :class="['counter-display', { pop: wssAnimating }]">
          <!-- The ?? operator ensures we show a dash perfectly if the value is null -->
          <span class="counter-value">{{ counter ?? '—' }}</span>
          <span class="counter-label">Global Clicks</span>
        </div>

        <!-- Master Mutation Button -->
        <Button
          id="increment-btn"
          class="increment-btn"
          :loading="incrementing"
          :disabled="incrementing || connectionStatus !== 'connected'"
          @click="increment"
          severity="primary"
          size="large"
          rounded
          icon="pi pi-bolt"
          :label="`Increment ${counterName || 'main'}`"
        />

        <Message v-if="error" severity="error" :closable="false" class="error-msg">
          {{ error }}
        </Message>

        <div class="meta-row">
          <i class="pi pi-clock meta-icon" />
          <span class="meta-text" v-if="counter === 0">Wait state - pending UPSERT</span>
          <span class="meta-text" v-else>DB Sync: {{ formatDate(updatedAt) }}</span>
        </div>
      </div>
    </template>
  </Card>
</template>

<style scoped>
/* NOTE: 'scoped' means any CSS class written here will NOT accidentally bleed over into App.vue */
.counter-card {
  width: 100%;
  max-width: 480px;
  background: rgba(24, 24, 27, 0.8) !important;
  border: 1px solid #3f3f46 !important;
  border-radius: 20px !important;
  backdrop-filter: blur(20px);
  box-shadow: 0 0 0 1px rgba(124,58,237,0.1), 0 25px 60px rgba(0,0,0,0.5) !important;
}

.counter-selector { width: 100%; margin-top: 10px; }
.counter-selector :deep(.p-inputtext) {
  width: 100%; background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); color: #fff; padding: 18px 14px 12px 14px;
}
.counter-selector :deep(.p-floatlabel label) { color: #a1a1aa; }

.spinner-wrap { display: flex; justify-content: center; padding: 40px 0; }
.counter-body { display: flex; flex-direction: column; align-items: center; gap: 28px; padding: 8px 0; border-top: 1px dashed rgba(255,255,255,0.1); margin-top: 25px; padding-top: 30px; }

.ws-status { display: flex; align-items: center; gap: 8px; background: rgba(0,0,0,0.3); padding: 6px 12px; border-radius: 99px; border: 1px solid rgba(255,255,255,0.05); }
.ws-status .dot { width: 8px; height: 8px; border-radius: 50%; }
.ws-status .dot.connected { background: #34d399; box-shadow: 0 0 8px rgba(52, 211, 153, 0.6); }
.ws-status .dot.connecting { background: #fbbf24; }
.ws-status .dot.disconnected { background: #ef4444; }
.ws-status .status-text { font-size: 0.75rem; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; color: #a1a1aa; }

.counter-display { display: flex; flex-direction: column; align-items: center; gap: 6px; transition: transform 0.1s ease; }
/* The .pop class defines our magical glow animation whenever a websocket message fires */
.counter-display.pop { animation: pop 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97); }
@keyframes pop {
  0%   { transform: scale(1); }
  10%  { transform: scale(1.3); filter: drop-shadow(0 0 8px #34d399); }
  50%  { transform: scale(0.95); }
  75%  { transform: scale(1.05); }
  100% { transform: scale(1); filter: drop-shadow(0 0 0 transparent); }
}

.counter-value {
  font-size: clamp(4rem, 12vw, 6rem); font-weight: 800; letter-spacing: -3px;
  background: linear-gradient(135deg, #7c3aed, #22d3ee); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; line-height: 1;
}
.counter-display.pop .counter-value {
  background: linear-gradient(135deg, #34d399, #10b981); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
}
.counter-label { font-size: 0.85rem; text-transform: uppercase; letter-spacing: 2px; color: #71717a; font-weight: 600; }

.increment-btn { width: 100%; max-width: 280px; height: 52px; font-size: 1rem !important; font-weight: 600 !important; background: linear-gradient(135deg, #6d28d9, #7c3aed) !important; border: none !important; box-shadow: 0 0 24px rgba(109,40,217,0.4) !important; transition: transform 0.15s ease, box-shadow 0.15s ease !important; }
.increment-btn:hover:not(:disabled) { transform: translateY(-2px) !important; box-shadow: 0 0 36px rgba(109,40,217,0.6) !important; }
.increment-btn:active:not(:disabled) { transform: translateY(0) !important; }
.increment-btn:disabled { opacity: 0.5 !important; cursor: not-allowed !important; }

.error-msg { width: 100%; }
.meta-row { display: flex; align-items: center; gap: 6px; color: #52525b; font-size: 0.8rem; }
.meta-icon { font-size: 0.75rem; }
</style>
