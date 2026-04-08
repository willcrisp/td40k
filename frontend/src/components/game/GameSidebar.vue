<script setup lang="ts">
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useWebSocketStore } from '@/stores/useWebSocketStore';
import { useRoomStore } from '@/stores/useRoomStore';

const ws = useWebSocketStore();
const room = useRoomStore();
const { connected } = storeToRefs(ws);
const { battleRound } = storeToRefs(room);

const activeNav = ref('battlelog');

const navItems = [
  { id: 'battlelog', icon: 'pi pi-video', label: 'BATTLELOG' },
  { id: 'mission', icon: 'pi pi-file', label: 'MISSION' },
  { id: 'stratagems', icon: 'pi pi-bolt', label: 'STRATAGEMS' },
  { id: 'units', icon: 'pi pi-chart-line', label: 'UNIT_DATA' },
];
</script>

<template>
  <aside class="game-sidebar">
    <!-- Sector Header -->
    <div class="sidebar-sector">
      <span class="sector-label font-mono">SECTOR</span>
      <span class="sector-number font-display">
        {{ String(battleRound).padStart(2, '0') }}
      </span>
    </div>

    <!-- Icon Nav -->
    <nav class="sidebar-nav">
      <button
        v-for="item in navItems"
        :key="item.id"
        class="sidebar-btn"
        :class="{ 'sidebar-btn--active': activeNav === item.id }"
        @click="activeNav = item.id"
      >
        <i :class="item.icon" class="sidebar-btn-icon"></i>
        <span class="sidebar-btn-label font-mono">{{ item.label }}</span>
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.game-sidebar {
  width: 80px;
  min-width: 80px;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: var(--surface-container-lowest);
  border-right: 1px solid var(--ghost-border);
  z-index: 10;
  padding-top: 0.5rem;
}

.sidebar-sector {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0.75rem 0;
  margin-bottom: 0.5rem;
  border-bottom: 1px solid var(--ghost-border);
  width: 100%;
}

.sector-label {
  font-size: 0.5625rem;
  letter-spacing: 0.12em;
  color: var(--on-surface-variant);
  text-transform: uppercase;
}

.sector-number {
  font-size: 1.25rem;
  color: var(--on-surface);
  font-weight: 700;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  gap: 2px;
}

.sidebar-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
  width: 100%;
  padding: 0.75rem 0.25rem;
  border: none;
  background: transparent;
  color: var(--on-surface-variant);
  cursor: pointer;
  transition: all 0.15s;
}

.sidebar-btn:hover {
  background: var(--surface-container-low);
  color: var(--on-surface);
}

.sidebar-btn--active {
  background: var(--primary-container);
  color: white;
}

.sidebar-btn-icon {
  font-size: 1.125rem;
}

.sidebar-btn-label {
  font-size: 0.5rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}
</style>
