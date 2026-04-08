<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useRoomStore } from '@/stores/useRoomStore';

const room = useRoomStore();
const { attackerId, defenderId, status } = storeToRefs(room);
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Attacker Slot -->
    <div
      class="deploy-card"
      :class="attackerId ? 'deploy-card--ready' : 'deploy-card--waiting'"
    >
      <div class="deploy-card-header">
        <span class="font-mono text-xs text-surface-variant">ATTACKER</span>
        <span
          class="deploy-status font-mono"
          :class="attackerId ? 'deploy-status--ready' : 'deploy-status--waiting'"
        >
          {{ attackerId ? 'CONNECTED' : 'WAITING' }}
        </span>
      </div>
      <div class="deploy-card-indicator">
        <i
          class="pi"
          :class="attackerId ? 'pi-check-circle' : 'pi-spin pi-spinner'"
          style="font-size: 1.25rem"
        ></i>
      </div>
    </div>

    <!-- Defender Slot -->
    <div
      class="deploy-card"
      :class="defenderId ? 'deploy-card--ready' : 'deploy-card--waiting'"
    >
      <div class="deploy-card-header">
        <span class="font-mono text-xs text-surface-variant">DEFENDER</span>
        <span
          class="deploy-status font-mono"
          :class="defenderId ? 'deploy-status--ready' : 'deploy-status--waiting'"
        >
          {{ defenderId ? 'CONNECTED' : 'WAITING' }}
        </span>
      </div>
      <div class="deploy-card-indicator">
        <i
          class="pi"
          :class="defenderId ? 'pi-check-circle' : 'pi-spin pi-spinner'"
          style="font-size: 1.25rem"
        ></i>
      </div>
    </div>
  </div>

  <p
    v-if="status === 'lobby'"
    class="font-mono text-surface-variant text-center mt-4"
    style="font-size: 0.6875rem; letter-spacing: 0.08em"
  >
    Share this page URL for others to join.
  </p>
</template>
