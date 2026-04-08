<template>
  <Dialog
    v-model:visible="visible"
    header="Import Army List"
    :modal="true"
    :closable="true"
    :style="{ width: '520px' }"
    @hide="onHide"
  >
    <div v-if="!importResult" class="import-form">
      <p class="hint">
        Paste your ListForge JSON export below. In ListForge, use
        <strong>Share → Copy JSON</strong> to get the export.
      </p>
      <Textarea
        v-model="rawJson"
        rows="10"
        placeholder='{ "roster": { ... } }'
        :disabled="loading"
        class="json-input"
      />
      <div v-if="parseError" class="parse-error">
        {{ parseError }}
      </div>
    </div>

    <div v-else class="import-result">
      <div class="result-summary">
        <Tag severity="success" :value="importResult.faction_name" />
        <span class="pts">{{ importResult.total_points }} pts</span>
      </div>

      <div class="result-section">
        <div class="section-label">
          Matched ({{ importResult.matched.length }})
        </div>
        <div
          v-for="u in importResult.matched"
          :key="u.datasheet_id + u.name"
          class="result-row matched"
        >
          <span class="unit-name">
            <span v-if="u.quantity > 1" class="qty">{{ u.quantity }}x</span>
            {{ u.name }}
          </span>
          <span class="unit-pts">{{ u.points }} pts</span>
        </div>
      </div>

      <div v-if="importResult.unmatched.length > 0" class="result-section">
        <div class="section-label warn">
          Unmatched ({{ importResult.unmatched.length }}) — Wahapedia sync
          required or names differ
        </div>
        <div
          v-for="name in importResult.unmatched"
          :key="name"
          class="result-row unmatched"
        >
          {{ name }}
        </div>
      </div>
    </div>

    <template #footer>
      <div class="footer-actions">
        <Button
          v-if="!importResult"
          label="Import"
          :loading="loading"
          :disabled="!rawJson.trim() || loading"
          @click="doImport"
        />
        <Button
          v-if="importResult"
          label="Done"
          @click="onDone"
        />
        <Button
          v-if="importResult"
          label="Re-import"
          severity="secondary"
          @click="reset"
        />
        <Button
          label="Cancel"
          severity="secondary"
          :disabled="loading"
          @click="onCancel"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import Textarea from 'primevue/textarea';
import Tag from 'primevue/tag';
import { useRosterStore } from '@/stores/useRosterStore';
import type { ImportRosterResponse } from '@/types';

const props = defineProps<{ roomId: string }>();
const emit = defineEmits<{ (e: 'close'): void }>();

const rosterStore = useRosterStore();

const visible = ref(true);
const rawJson = ref('');
const parseError = ref('');
const loading = ref(false);
const importResult = ref<ImportRosterResponse | null>(null);

async function doImport() {
  parseError.value = '';
  let parsed: unknown;
  try {
    parsed = JSON.parse(rawJson.value);
  } catch {
    parseError.value = 'Invalid JSON — make sure you copied the full export.';
    return;
  }

  loading.value = true;
  try {
    const result = await rosterStore.importRoster(props.roomId, parsed);
    importResult.value = result;
  } catch (err: unknown) {
    const msg =
      err instanceof Error ? err.message : 'Import failed. Check the console.';
    parseError.value = msg;
  } finally {
    loading.value = false;
  }
}

function reset() {
  importResult.value = null;
  rawJson.value = '';
  parseError.value = '';
}

function onDone() {
  visible.value = false;
  emit('close');
}

function onCancel() {
  visible.value = false;
  emit('close');
}

function onHide() {
  emit('close');
}
</script>

<style scoped>
.hint {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.65);
  margin-bottom: 0.75rem;
}

.json-input {
  width: 100%;
  font-family: monospace;
  font-size: 0.75rem;
}

.parse-error {
  margin-top: 0.5rem;
  color: #f97316;
  font-size: 0.875rem;
}

.import-result {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.result-summary {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.pts {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.result-section {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.section-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 0.25rem;
}

.section-label.warn {
  color: #f97316;
}

.result-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
  padding: 0.25rem 0.5rem;
  border-radius: 3px;
}

.result-row.matched {
  background: rgba(34, 197, 94, 0.08);
}

.result-row.unmatched {
  background: rgba(249, 115, 22, 0.08);
  color: #f97316;
}

.unit-name {
  display: flex;
  gap: 0.375rem;
  align-items: center;
}

.qty {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
}

.unit-pts {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
}

.footer-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}
</style>
