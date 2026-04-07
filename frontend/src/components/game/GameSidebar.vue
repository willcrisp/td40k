<script setup lang="ts">
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useWebSocketStore } from '@/stores/useWebSocketStore';

const ws = useWebSocketStore();
const { connected } = storeToRefs(ws);
const showHistory = ref(false);

const navItems = [
  { icon: 'pi pi-wifi', label: 'CONNECTION', active: true },
  { icon: 'pi pi-history', label: 'HISTORY', active: false, onClick: () => (showHistory.value = !showHistory.value) },
  { icon: 'pi pi-box', label: 'LOGISTICS', active: false },
  { icon: 'pi pi-bolt', label: 'ABILITIES', active: false },
  { icon: 'pi pi-cog', label: 'SETTINGS', active: false },
];

const logEntries = [
  { time: '14:32:01', text: 'Attacker moved units' },
  { time: '14:34:45', text: 'Unit-043 deployed smoke' },
  { time: '14:35:12', text: 'Terrain: Ruins (Obstacle)' },
];
</script>

<template>
  <aside class="game-sidebar panel-lowest riveted">
    <!-- Sector Header -->
    <div class="sidebar-header panel-low">
      <div class="flex items-center justify-between mb-1">
        <h2 class="font-display text-[11px] text-white opacity-80 uppercase tracking-widest">Active Game</h2>
        <div
          class="status-dot"
          :class="connected ? 'status-dot--online' : 'status-dot--offline'"
        />
      </div>
      <p class="text-[9px] font-mono" :class="connected ? 'text-primary' : 'text-on-surface-variant'">
        {{ connected ? 'CONNECTION STABLE' : 'SIGNAL LOST' }}
      </p>
    </div>

    <!-- Navigation -->
    <nav class="sidebar-nav">
      <button
        v-for="item in navItems"
        :key="item.label"
        class="nav-item"
        :class="{ 'nav-item--active': item.active || (item.label === 'HISTORY' && showHistory) }"
        @click="item.onClick ? item.onClick() : null"
      >
        <i :class="item.icon" class="text-xs" />
        <span class="font-display tracking-widest">{{ item.label }}</span>
      </button>
    </nav>

    <!-- Tactical Log Overlay/Panel -->
    <div
      v-if="showHistory"
      class="tactical-log-panel panel-low border-t border-outline-variant"
    >
      <div class="flex items-center justify-between mb-4">
        <span class="font-display text-[10px] text-white opacity-60">TACTICAL HISTORY</span>
        <button @click="showHistory = false" class="text-[10px] uppercase font-mono opacity-40 hover:opacity-100 cursor-pointer">
          CLOSE
        </button>
      </div>
      <div class="log-entries custom-scrollbar">
        <div
          v-for="entry in logEntries"
          :key="entry.time"
          class="log-entry"
        >
          <div class="flex items-baseline gap-2">
            <span class="log-time font-mono text-[8px] text-surface-variant">[{{ entry.time }}]</span>
            <span class="log-text font-mono text-[10px] text-on-surface">{{ entry.text }}</span>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.game-sidebar {
  width: 260px;
  min-width: 260px;
  height: 100%;
  display: flex;
  flex-direction: column;
  z-index: 10;
}

.sidebar-sector {
  padding: 24px 20px;
  margin-bottom: 2px;
}

.sector-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 0;
}

.status-dot--online {
  background: var(--tertiary);
  box-shadow: 0 0 10px var(--tertiary-glow);
}

.status-dot--offline {
  background: var(--primary);
  box-shadow: 0 0 10px var(--primary-container);
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  padding: 12px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  border: none;
  background: transparent;
  color: var(--on-surface-variant);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s steps(2);
  text-align: left;
}

.nav-item:hover {
  background: var(--surface-container-low);
  color: var(--on-surface);
}

.nav-item--active {
  background: var(--surface-container-high);
  color: var(--tertiary);
  border-left: 4px solid var(--tertiary);
  padding-left: 16px;
}

.tactical-log-panel {
  padding: 1rem;
  margin-top: auto;
  max-height: 200px;
  display: flex;
  flex-direction: column;
}

.log-entries {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
  padding-right: 4px;
}

.log-entry {
  line-height: 1.2;
}

.log-text {
  text-transform: none; /* Already distilled to plain English */
}

/* Scrollbar refinement */
.custom-scrollbar::-webkit-scrollbar {
  width: 2px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--outline-variant);
}
</style>
