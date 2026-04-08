<template>
  <div class="unit-placer">
    <!-- Tab bar -->
    <div v-if="!isPlacing" class="tab-bar">
      <button
        :class="['tab-btn', { active: activeTab === 'datasheet' }]"
        @click="activeTab = 'datasheet'"
      >
        Datasheets
      </button>
      <button
        :class="['tab-btn', { active: activeTab === 'roster' }]"
        @click="activeTab = 'roster'"
      >
        Roster
        <span v-if="rosterStore.roster.length > 0" class="roster-count">
          {{ rosterStore.roster.length }}
        </span>
      </button>
    </div>

    <!-- Datasheet tab -->
    <div v-if="!isPlacing && activeTab === 'datasheet'" class="placer-form">
      <div class="form-group">
        <label>Datasheet</label>
        <Dropdown
          v-model="selectedDatasheet"
          :options="datasheets"
          option-label="name"
          option-value="id"
          placeholder="Select unit type..."
          @change="onDatasheetChange"
          filter
          :loading="loadingDatasheets"
        />
      </div>

      <div v-if="selectedDatasheet && models.length > 0" class="form-group">
        <label>Model</label>
        <Dropdown
          v-model="selectedModel"
          :options="models"
          option-label="name"
          placeholder="Select model..."
          :loading="loadingModels"
        />
      </div>

      <div v-if="selectedModel" class="form-group">
        <label>Model Count</label>
        <InputNumber
          v-model="modelCount"
          :min="1"
          :max="20"
          placeholder="Number of models"
        />
      </div>

      <div v-if="selectedModel && selectedDatasheet" class="unit-preview">
        <div class="preview-header">Unit Preview</div>
        <div v-if="selectedDatasheet" class="stats-grid">
          <div class="stat">
            <span class="label">Movement:</span>
            <span class="value">{{ selectedModel?.m || '-' }}"</span>
          </div>
          <div class="stat">
            <span class="label">Toughness:</span>
            <span class="value">{{ selectedModel?.t || '-' }}</span>
          </div>
          <div class="stat">
            <span class="label">Save:</span>
            <span class="value">{{ selectedModel?.sv || '-' }}</span>
          </div>
          <div class="stat">
            <span class="label">Wounds:</span>
            <span class="value">{{ selectedModel?.w || '-' }}</span>
          </div>
        </div>

        <div class="footprint-preview">
          <svg viewBox="0 0 100 100" class="footprint-svg">
            <rect
              x="0"
              y="0"
              width="100"
              height="100"
              fill="rgba(128,128,128,0.1)"
              stroke="rgba(128,128,128,0.3)"
              stroke-width="1"
            />
            <component
              :is="footprintComponent"
              :base-size="selectedModel?.base_size"
            />
          </svg>
          <div class="footprint-label">
            Base: {{ selectedModel?.base_size_descr || selectedModel?.base_size || 'N/A' }}
          </div>
        </div>
      </div>

      <div v-if="selectedModel" class="form-actions">
        <Button
          label="Place Unit"
          @click="startPlacing"
          :disabled="!selectedModel || modelCount < 1"
          class="p-button-success"
        />
      </div>
    </div>

    <!-- Roster tab -->
    <div v-if="!isPlacing && activeTab === 'roster'" class="roster-panel">
      <Button
        label="Import List"
        severity="secondary"
        size="small"
        @click="showImportModal = true"
        class="import-btn"
      />

      <div v-if="rosterStore.roster.length === 0" class="roster-empty">
        No roster imported yet. Click "Import List" to paste your ListForge
        export.
      </div>

      <div v-else class="roster-list">
        <div
          v-for="entry in rosterStore.roster"
          :key="entry.id || entry.datasheet_id"
          class="roster-entry"
          @click="placeFromRoster(entry)"
        >
          <div class="entry-name">
            <span v-if="entry.quantity > 1" class="entry-qty">
              {{ entry.quantity }}x
            </span>
            {{ entry.model_name }}
          </div>
          <span class="entry-pts">{{ entry.points }} pts</span>
        </div>
      </div>

      <div v-if="rosterStore.roster.length > 0" class="roster-actions">
        <Button
          label="Clear Roster"
          severity="danger"
          size="small"
          @click="clearRoster"
        />
      </div>
    </div>

    <!-- Placing instruction -->
    <div v-if="isPlacing" class="placing-instruction">
      <div class="instruction-icon">🎯</div>
      <div class="instruction-text">
        Click on the board to place<br />
        <strong>{{ modelCount }}x {{ selectedModel?.name }}</strong>
      </div>
      <Button
        label="Cancel"
        @click="cancelPlacing"
        class="p-button-secondary"
      />
    </div>

    <RosterImportModal
      v-if="showImportModal"
      :room-id="roomId"
      @close="showImportModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';
import Button from 'primevue/button';
import { useUnitStore } from '@/stores/useUnitStore';
import { useRosterStore } from '@/stores/useRosterStore';
import {
  apiGetDatasheets,
  apiGetDatasheetModels,
  type WhDatasheet,
  type WhDatasheetModel,
} from '@/lib/api';
import type { RosterEntry } from '@/types';
import RosterImportModal from './RosterImportModal.vue';

const props = defineProps<{ roomId: string }>();

const unitStore = useUnitStore();
const rosterStore = useRosterStore();

const datasheets = ref<WhDatasheet[]>([]);
const models = ref<WhDatasheetModel[]>([]);
const selectedDatasheet = ref<string | null>(null);
const selectedModel = ref<WhDatasheetModel | null>(null);
const modelCount = ref(1);

const loadingDatasheets = ref(false);
const loadingModels = ref(false);

const activeTab = ref<'datasheet' | 'roster'>('datasheet');
const showImportModal = ref(false);

const isPlacing = computed(() => unitStore.placingUnitType !== null);

const footprintComponent = computed(() => {
  if (!selectedModel.value) return 'div';
  const baseSize = selectedModel.value.base_size?.toLowerCase() || '';

  if (
    baseSize === '25mm' ||
    baseSize === '32mm' ||
    baseSize === '40mm' ||
    baseSize === '50mm' ||
    baseSize === '60mm'
  ) {
    return FootprintCircle;
  }
  if (baseSize.includes('x')) {
    return FootprintOval;
  }
  if (baseSize.includes('hull')) {
    return FootprintHull;
  }
  return FootprintCircle;
});

onMounted(async () => {
  loadingDatasheets.value = true;
  try {
    const { data } = await apiGetDatasheets();
    datasheets.value = data;
  } catch (err) {
    console.error('Failed to load datasheets:', err);
  } finally {
    loadingDatasheets.value = false;
  }

  // Load persisted roster for this room
  if (props.roomId) {
    await rosterStore.loadRoster(props.roomId).catch(() => {
      // Non-fatal — roster may be empty
    });
  }
});

async function onDatasheetChange() {
  if (!selectedDatasheet.value) {
    models.value = [];
    selectedModel.value = null;
    return;
  }

  loadingModels.value = true;
  try {
    const { data } = await apiGetDatasheetModels(selectedDatasheet.value);
    models.value = data;
    selectedModel.value = data.length > 0 ? data[0] : null;
  } catch (err) {
    console.error('Failed to load models:', err);
    models.value = [];
  } finally {
    loadingModels.value = false;
  }
}

function startPlacing() {
  if (!selectedDatasheet.value || !selectedModel.value) return;

  const ds = datasheets.value.find((d) => d.id === selectedDatasheet.value);
  if (!ds) return;

  unitStore.startPlacingUnit(
    selectedDatasheet.value,
    selectedModel.value.name,
    modelCount.value,
    ds.faction_id
  );
}

function cancelPlacing() {
  unitStore.cancelPlacing();
}

function placeFromRoster(entry: RosterEntry) {
  unitStore.startPlacingUnit(
    entry.datasheet_id,
    entry.model_name,
    entry.quantity,
    entry.faction_id
  );
}

async function clearRoster() {
  if (!props.roomId) return;
  await rosterStore.clearRoster(props.roomId);
}
</script>

<style scoped>
.unit-placer {
  padding: 1rem;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.tab-bar {
  display: flex;
  gap: 0.25rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding-bottom: 0.5rem;
}

.tab-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8125rem;
  font-weight: 500;
  padding: 0.25rem 0.625rem;
  border-radius: 3px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.tab-btn:hover {
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.05);
}

.tab-btn.active {
  color: rgba(255, 255, 255, 0.95);
  background: rgba(255, 255, 255, 0.1);
}

.roster-count {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 10px;
  padding: 0 0.375rem;
  font-size: 0.7rem;
}

.placer-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.unit-preview {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 1rem;
  gap: 1rem;
  display: flex;
  flex-direction: column;
}

.preview-header {
  font-size: 0.875rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
  font-size: 0.875rem;
}

.stat {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 2px;
}

.stat .label {
  color: rgba(255, 255, 255, 0.6);
  font-weight: 500;
}

.stat .value {
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
}

.footprint-preview {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  align-items: center;
}

.footprint-svg {
  width: 80px;
  height: 80px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}

.footprint-label {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.6);
}

.form-actions {
  display: flex;
  gap: 0.5rem;
}

.roster-panel {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.import-btn {
  align-self: flex-start;
}

.roster-empty {
  font-size: 0.8125rem;
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
  padding: 1.5rem 0.5rem;
  line-height: 1.5;
}

.roster-list {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.roster-entry {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.75rem;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  cursor: pointer;
  font-size: 0.875rem;
}

.roster-entry:hover {
  background: rgba(255, 255, 255, 0.09);
  border-color: rgba(255, 255, 255, 0.18);
}

.entry-name {
  display: flex;
  gap: 0.375rem;
  align-items: center;
  color: rgba(255, 255, 255, 0.85);
}

.entry-qty {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.45);
}

.entry-pts {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.4);
}

.roster-actions {
  display: flex;
  justify-content: flex-end;
}

.placing-instruction {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 2rem 1rem;
  background: linear-gradient(
    135deg,
    rgba(33, 150, 243, 0.1),
    rgba(76, 175, 80, 0.1)
  );
  border: 2px dashed rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.instruction-icon {
  font-size: 2rem;
}

.instruction-text {
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
}

.instruction-text strong {
  color: rgba(255, 255, 255, 1);
  font-weight: 600;
}
</style>

<!-- Footprint shape components -->
<script lang="ts">
const FootprintCircle = {
  template:
    '<circle cx="50" cy="50" r="30" fill="rgba(100, 200, 255, 0.3)" stroke="rgba(100, 200, 255, 0.6)" stroke-width="2" />',
};

const FootprintOval = {
  template:
    '<ellipse cx="50" cy="50" rx="40" ry="25" fill="rgba(255, 150, 100, 0.3)" stroke="rgba(255, 150, 100, 0.6)" stroke-width="2" />',
};

const FootprintHull = {
  template:
    '<rect x="20" y="30" width="60" height="40" fill="rgba(200, 100, 100, 0.3)" stroke="rgba(200, 100, 100, 0.6)" stroke-width="2" rx="2" />',
};
</script>
