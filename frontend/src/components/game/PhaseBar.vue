<script setup lang="ts">
import { computed } from 'vue';
import { storeToRefs } from 'pinia';
import { PHASES, PHASE_LABELS } from '@/types';
import { useRoomStore } from '@/stores/useRoomStore';

const room = useRoomStore();
const { currentPhase } = storeToRefs(room);

const activeIndex = computed(() => PHASES.indexOf(currentPhase.value));

// Node width is 180px, spacing is 40px (20% of 200px total step)
const slideTransform = computed(() => {
  const step = 220; // 180 width + 40 gap
  const centerOffset = activeIndex.value * step;
  // Center is 50% of parent minus half the node width (90) minus the active index offset
  return `translateX(calc(50% - ${centerOffset}px - 90px))`;
});

function phaseLabel(phase: string) {
  return PHASE_LABELS[phase as keyof typeof PHASE_LABELS].replace(' Phase', '').toUpperCase();
}
</script>

<template>
  <div class="phase-bar px-4 relative overflow-hidden panel-low border-t border-b border-outline-variant">
    <!-- Edge Shadows -->
    <div class="phase-mask--left"></div>
    <div class="phase-mask--right"></div>

    <!-- Sliding Container -->
    <div 
      class="phase-slider flex items-center h-full transition-transform duration-700 ease-out fill-mode-forwards"
      :style="{ transform: slideTransform }"
    >
      <div
        v-for="(phase, idx) in PHASES"
        :key="phase"
        class="phase-node flex-shrink-0 flex items-center justify-center"
        style="width: 180px; margin: 0 20px;"
      >
        <div
          class="phase-content flex flex-col items-center transition-all duration-500"
          :class="{
            'phase-content--active': currentPhase === phase,
            'phase-content--past': idx < activeIndex,
            'phase-content--future': idx > activeIndex
          }"
        >
          <span class="font-display text-[10px] tracking-[0.3em] uppercase">
            {{ phaseLabel(phase) }}
          </span>
          <div v-if="currentPhase === phase" class="phase-indicator mt-1"></div>
        </div>
      </div>
    </div>

    <!-- Navigation Overlays -->
    <div class="phase-nav-container absolute inset-0 flex items-center justify-between px-4 pointer-events-none">
      <button 
        class="nav-btn pointer-events-auto"
        @click="room.prevPhase()"
      >
        <i class="pi pi-chevron-left text-[10px]" />
      </button>
      <button 
        class="nav-btn pointer-events-auto"
        @click="room.nextPhase()"
      >
        <i class="pi pi-chevron-right text-[10px]" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.phase-bar {
  height: 48px;
  width: 100%;
}

.phase-content {
  opacity: 0.2;
  transform: scale(0.9);
}

.phase-content--active {
  opacity: 1;
  transform: scale(1.1);
  color: var(--primary);
}

.phase-content--past {
  opacity: 0.5;
  color: var(--on-surface-variant);
}

.phase-indicator {
  width: 24px;
  height: 2px;
  background-color: var(--primary);
  box-shadow: 0 0 8px var(--primary);
}

.phase-mask--left, .phase-mask--right {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 15%;
  z-index: 10;
  pointer-events: none;
}

.phase-mask--left {
  left: 0;
  background: linear-gradient(to right, var(--surface-container-low), transparent);
}

.phase-mask--right {
  right: 0;
  background: linear-gradient(to left, var(--surface-container-low), transparent);
}

.nav-btn {
  background: none;
  border: none;
  color: var(--on-surface-variant);
  opacity: 0.4;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-btn:hover {
  opacity: 1;
  color: var(--on-surface);
  transform: scale(1.2);
}
</style>
