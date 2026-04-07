<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useWebSocketStore } from '@/stores/useWebSocketStore';
import { useRoomStore } from '@/stores/useRoomStore';
import { useBoardStore } from '@/stores/useBoardStore';
import PhaseBar from './PhaseBar.vue';

const router = useRouter();
const ws = useWebSocketStore();
const room = useRoomStore();
const board = useBoardStore();
const { status, winner, battleRound } = storeToRefs(room);

function goHome() {
  ws.disconnect();
  router.push('/');
}
</script>

<template>
  <header class="game-header panel-high border-b border-outline-variant">
    <!-- Left: Branding -->
    <div class="header-brand" @click="goHome">
      <div class="flex items-center gap-3">
        <div class="w-1 h-8 bg-primary"></div>
        <div class="flex flex-col">
          <span class="brand-title font-display text-white">GAME TRACKER V1.0</span>
          <span class="text-[9px] font-mono text-surface-variant">Session ID: 772-91-X</span>
        </div>
      </div>
    </div>

    <!-- Center: Primary Game Status -->
    <div class="header-center justify-center">
      <div class="flex items-center gap-12">
        <div class="flex flex-col items-center">
          <span class="text-[9px] font-mono text-tertiary uppercase opacity-60">Battle Round</span>
          <span class="text-xl font-display text-white">{{ battleRound }} / 5</span>
        </div>
        
        <div class="w-[1px] h-8 bg-outline-variant opacity-20"></div>

        <div class="flex flex-col items-center">
          <span class="text-[9px] font-mono text-tertiary uppercase opacity-60">Active Player</span>
          <span 
            class="text-xl font-display"
            :class="room.active_player === 'attacker' ? 'text-primary' : 'text-secondary'"
          >
            {{ room.active_player?.toUpperCase() }}
          </span>
        </div>
      </div>
    </div>

    <!-- Right: Icons -->
    <div class="header-actions">
      <button 
        class="header-icon-btn btn-secondary-tactical" 
        title="Rotate Board (90°)"
        @click="board.toggleRotation()"
      >
        <i class="pi pi-refresh" />
      </button>
      <button class="header-icon-btn btn-secondary-tactical" title="Settings">
        <i class="pi pi-cog" />
      </button>
      <button class="header-icon-btn btn-secondary-tactical" title="Fullscreen">
        <i class="pi pi-expand" />
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

  <!-- Phase Bar (below header) -->
  <div class="phase-bar-row panel-low border-b border-ghost-border">
    <PhaseBar />
  </div>
</template>

<style scoped>
.game-header {
  display: flex;
  align-items: center;
  height: 60px;
  min-height: 60px;
  padding: 0 20px;
  gap: 24px;
  z-index: 20;
  position: relative;
}

.header-brand {
  cursor: pointer;
  padding-right: 24px;
  border-right: 1px solid var(--ghost-border);
}

.brand-title {
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0.1em;
}

.header-center {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 32px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  font-size: 14px;
}

.game-over-banner {
  position: absolute;
  top: 70px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6);
}

.phase-bar-row {
  height: 48px;
  min-height: 48px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  z-index: 10;
}
</style>
