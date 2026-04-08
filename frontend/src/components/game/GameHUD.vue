<script setup lang="ts">
import { computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useWebSocketStore } from '@/stores/useWebSocketStore';
import { useRoomStore } from '@/stores/useRoomStore';
import { usePlayerStore } from '@/stores/usePlayerStore';
import { PHASES, PHASE_LABELS } from '@/types';

const router = useRouter();
const ws = useWebSocketStore();
const room = useRoomStore();
const playerStore = usePlayerStore();
const { status, winner, battleRound, currentPhase } = storeToRefs(room);

const activePhaseLabel = computed(() =>
  PHASE_LABELS[currentPhase.value].replace(' Phase', '').toUpperCase()
);

const nextPhaseLabel = computed(() => {
  const idx = PHASES.indexOf(currentPhase.value);
  if (idx < PHASES.length - 1) {
    return PHASE_LABELS[PHASES[idx + 1]]
      .replace(' Phase', '')
      .toUpperCase();
  }
  return 'END';
});

function goHome() {
  ws.disconnect();
  router.push('/');
}

function handleLogout() {
  playerStore.logout();
  router.push('/auth');
}
</script>

<template>
  <header class="game-header">
    <!-- Left: Branding + Round -->
    <div class="header-left">
      <span
        class="header-brand font-mono"
        style="cursor: pointer"
        @click="goHome"
      >TACTICAL TERMINAL</span>

      <div class="round-badge">
        <span class="round-badge-label font-mono">BATTLE ROUND</span>
        <span class="round-badge-number font-display">
          {{ String(battleRound).padStart(2, '0') }}
        </span>
      </div>
    </div>

    <!-- Center: Active Phase + Next -->
    <div class="header-phase">
      <div class="phase-active">
        <span class="phase-meta font-mono">ACTIVE PHASE</span>
        <span class="phase-name font-display">{{ activePhaseLabel }}</span>
      </div>
      <div class="phase-next">
        <span class="phase-meta font-mono">NEXT</span>
        <span class="phase-next-name font-mono">{{ nextPhaseLabel }}</span>
      </div>
    </div>

    <!-- Right: Tabs + Avatar -->
    <div class="header-right">
      <div class="header-tabs">
        <button class="header-tab header-tab--active font-mono">
          <i class="pi pi-exclamation-triangle" style="font-size: 0.625rem"></i>
          PRIMARY
        </button>
        <button class="header-tab font-mono">
          <i class="pi pi-check-square" style="font-size: 0.625rem"></i>
          SECONDARY
        </button>
        <button class="header-tab font-mono">
          <i class="pi pi-book" style="font-size: 0.625rem"></i>
          RULES
        </button>
      </div>
      <button class="header-avatar" @click="handleLogout">
        <i class="pi pi-user"></i>
      </button>
    </div>

    <!-- Game over banner -->
    <div
      v-if="status === 'finished'"
      class="game-over-banner bg-anodized-primary riveted"
    >
      <div class="flex items-center gap-4 px-6 py-2">
        <div class="status-seal">GG</div>
        <div class="flex flex-col">
          <span class="font-display text-lg tracking-widest text-white">GAME OVER</span>
          <span v-if="winner" class="text-xs font-mono text-white opacity-80">
            {{ winner.toUpperCase() }} WON
          </span>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.game-header {
  display: flex;
  align-items: center;
  height: 52px;
  min-height: 52px;
  padding: 0 1rem;
  background-color: var(--surface-container-lowest);
  border-top: 2px solid var(--primary-container);
  border-bottom: 1px solid var(--ghost-border);
  z-index: 20;
  position: relative;
  gap: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding-right: 1.5rem;
  border-right: 1px solid var(--ghost-border);
}

.round-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background-color: var(--primary-container);
  padding: 0.25rem 0.625rem;
}

.round-badge-label {
  font-size: 0.5625rem;
  letter-spacing: 0.1em;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
}

.round-badge-number {
  font-size: 1.125rem;
  color: white;
  font-weight: 700;
}

.header-phase {
  display: flex;
  align-items: center;
  padding: 0 1.5rem;
  gap: 1.5rem;
  border-right: 1px solid var(--ghost-border);
}

.phase-active {
  display: flex;
  flex-direction: column;
}

.phase-meta {
  font-size: 0.5625rem;
  letter-spacing: 0.08em;
  color: var(--on-surface-variant);
  text-transform: uppercase;
}

.phase-name {
  font-size: 1rem;
  color: var(--primary);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.phase-next {
  display: flex;
  flex-direction: column;
}

.phase-next-name {
  font-size: 0.8125rem;
  color: var(--on-surface-variant);
  text-transform: uppercase;
}

.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-tabs {
  display: flex;
  gap: 2px;
}

.header-tab {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  font-size: 0.625rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  background: var(--surface-container-high);
  border: none;
  color: var(--on-surface-variant);
  cursor: pointer;
  transition: all 0.15s;
}

.header-tab:hover {
  background: var(--surface-container-highest);
  color: var(--on-surface);
}

.header-tab--active {
  background: var(--surface-container-highest);
  color: var(--on-surface);
}

.game-over-banner {
  position: absolute;
  top: 60px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6);
}
</style>
