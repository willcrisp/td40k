<template>
  <div v-if="unit" class="unit-inspector">
    <div class="inspector-header">
      <div class="unit-title">
        <div class="unit-name">{{ unit.name_on_board || unit.model_name }}</div>
        <div class="unit-count">{{ unit.model_count }}x {{ unit.model_name }}</div>
      </div>
      <Button
        icon="pi pi-times"
        class="p-button-rounded p-button-text p-button-sm"
        @click="closeInspector"
      />
    </div>

    <div class="inspector-content">
      <div class="section status-section">
        <div class="section-title">Status</div>
        <div class="status-display">
          <Tag
            :value="unit.status"
            :severity="statusSeverity"
            class="status-tag"
          />
          <span v-if="unit.status === 'alive'" class="status-text">
            {{ unit.wounds }} / {{ maxWounds }} wounds
          </span>
          <span v-else class="status-text">
            {{ statusLabel }}
          </span>
        </div>
      </div>

      <div v-if="unit.status === 'alive'" class="section wounds-section">
        <div class="section-title">Manage Wounds</div>
        <div class="wounds-controls">
          <Button
            label="-"
            @click="decreaseWounds"
            :disabled="unit.wounds <= 0"
            class="p-button-sm"
          />
          <div class="wound-count">{{ unit.wounds }}</div>
          <Button
            label="+"
            @click="increaseWounds"
            :disabled="unit.wounds >= maxWounds"
            class="p-button-sm"
          />
        </div>
      </div>

      <div class="section status-actions">
        <div class="section-title">Actions</div>
        <div class="button-group">
          <Button
            v-if="unit.status === 'alive'"
            label="Send to Reserves"
            @click="sendToReserves"
            severity="warning"
            class="p-button-sm"
          />
          <Button
            v-if="unit.status === 'alive'"
            label="Mark Dead"
            @click="markDead"
            severity="danger"
            class="p-button-sm"
          />
          <Button
            v-if="unit.status === 'in_reserves'"
            label="Deploy"
            @click="deploy"
            severity="success"
            class="p-button-sm"
          />
          <Button
            v-if="unit.status === 'dead'"
            label="Revive"
            @click="revive"
            severity="warning"
            class="p-button-sm"
          />
          <Button
            label="Remove from Board"
            @click="removeUnit"
            severity="danger"
            class="p-button-sm"
          />
        </div>
      </div>
    </div>
  </div>
  <div v-else class="inspector-empty">
    <div class="empty-state">
      <div class="empty-icon">📋</div>
      <div class="empty-text">Select a unit to inspect</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import {
  apiWoundUnit,
  apiUpdateUnitStatus,
  apiRemoveUnit,
} from '@/lib/api';
import { useUnitStore } from '@/stores/useUnitStore';
import { useRoomStore } from '@/stores/useRoomStore';
import type { GameUnit } from '@/types';

const unitStore = useUnitStore();
const roomStore = useRoomStore();

const unit = computed(() => unitStore.selectedUnit);

const maxWounds = computed(() => {
  // TODO: Get max wounds from unit stats once we fetch them
  // For now, estimate based on model count
  return (unit.value?.model_count || 1) * 2; // rough estimate
});

const statusLabel = computed(() => {
  switch (unit.value?.status) {
    case 'in_reserves':
      return 'In Reserves';
    case 'dead':
      return 'Removed from Battle';
    case 'alive':
    default:
      return 'Active';
  }
});

const statusSeverity = computed(() => {
  switch (unit.value?.status) {
    case 'alive':
      return 'success';
    case 'in_reserves':
      return 'warning';
    case 'dead':
      return 'danger';
    default:
      return 'info';
  }
});

function closeInspector() {
  unitStore.selectUnit(null);
}

async function increaseWounds() {
  if (!unit.value || !roomStore.roomId) return;
  try {
    const updated = await apiWoundUnit(
      roomStore.roomId,
      unit.value.id,
      1
    );
    unitStore.updateUnit(unit.value.id, { wounds: updated.data.wounds });
  } catch (err) {
    console.error('Failed to apply wound:', err);
  }
}

function decreaseWounds() {
  if (!unit.value || unit.value.wounds <= 0) return;
  // Note: We could implement a "heal" endpoint if needed
  // For now, wounds only increase
}

async function sendToReserves() {
  if (!unit.value || !roomStore.roomId) return;
  try {
    const updated = await apiUpdateUnitStatus(
      roomStore.roomId,
      unit.value.id,
      'in_reserves'
    );
    unitStore.updateUnit(unit.value.id, {
      status: updated.data.status,
    });
  } catch (err) {
    console.error('Failed to send to reserves:', err);
  }
}

async function markDead() {
  if (!unit.value || !roomStore.roomId) return;
  try {
    const updated = await apiUpdateUnitStatus(
      roomStore.roomId,
      unit.value.id,
      'dead'
    );
    unitStore.updateUnit(unit.value.id, { status: updated.data.status });
  } catch (err) {
    console.error('Failed to mark dead:', err);
  }
}

async function deploy() {
  if (!unit.value || !roomStore.roomId) return;
  try {
    const updated = await apiUpdateUnitStatus(
      roomStore.roomId,
      unit.value.id,
      'alive'
    );
    unitStore.updateUnit(unit.value.id, { status: updated.data.status });
  } catch (err) {
    console.error('Failed to deploy:', err);
  }
}

async function revive() {
  if (!unit.value || !roomStore.roomId) return;
  try {
    const updated = await apiUpdateUnitStatus(
      roomStore.roomId,
      unit.value.id,
      'alive'
    );
    unitStore.updateUnit(unit.value.id, { status: updated.data.status });
  } catch (err) {
    console.error('Failed to revive:', err);
  }
}

async function removeUnit() {
  if (!unit.value || !roomStore.roomId) return;
  try {
    await apiRemoveUnit(roomStore.roomId, unit.value.id);
    unitStore.removeUnit(unit.value.id);
  } catch (err) {
    console.error('Failed to remove unit:', err);
  }
}
</script>

<style scoped>
.unit-inspector {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.inspector-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: rgba(33, 150, 243, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.unit-title {
  flex: 1;
}

.unit-name {
  font-size: 1rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 1);
}

.unit-count {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.6);
  margin-top: 0.25rem;
}

.inspector-content {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  overflow-y: auto;
}

.section {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.02);
}

.section-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
}

.status-section {
  background: rgba(33, 150, 243, 0.1);
  border-color: rgba(33, 150, 243, 0.3);
}

.status-display {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.status-tag {
  flex-shrink: 0;
}

.status-text {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.8);
}

.wounds-section {
  background: rgba(255, 152, 0, 0.1);
  border-color: rgba(255, 152, 0, 0.3);
}

.wounds-controls {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.wound-count {
  flex: 1;
  text-align: center;
  font-size: 1.25rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  padding: 0.5rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 2px;
}

.status-actions {
  background: rgba(76, 175, 80, 0.1);
  border-color: rgba(76, 175, 80, 0.3);
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.button-group :deep(.p-button) {
  width: 100%;
  font-size: 0.75rem;
}

.inspector-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  min-height: 200px;
  background: rgba(0, 0, 0, 0.1);
  border: 1px dashed rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  text-align: center;
}

.empty-icon {
  font-size: 2rem;
  opacity: 0.5;
}

.empty-text {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.875rem;
}
</style>
