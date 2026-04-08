import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { RosterEntry, ImportRosterResponse } from '@/types';
import {
  apiImportRoster,
  apiGetRoster,
  apiClearRoster,
} from '@/lib/api';

export const useRosterStore = defineStore('roster', () => {
  const roster = ref<RosterEntry[]>([]);
  const lastImport = ref<ImportRosterResponse | null>(null);
  const loading = ref(false);

  async function importRoster(
    roomId: string,
    listforgeJson: unknown
  ): Promise<ImportRosterResponse> {
    loading.value = true;
    try {
      const { data } = await apiImportRoster(roomId, listforgeJson);
      lastImport.value = data;
      roster.value = data.matched.map((m) => ({
        id: '',
        room_id: roomId,
        player_id: '',
        datasheet_id: m.datasheet_id,
        model_name: m.name,
        quantity: m.quantity,
        faction_id: m.faction_id,
        points: m.points,
        created_at: '',
      }));
      return data;
    } finally {
      loading.value = false;
    }
  }

  async function loadRoster(roomId: string): Promise<void> {
    loading.value = true;
    try {
      const { data } = await apiGetRoster(roomId);
      roster.value = data ?? [];
    } finally {
      loading.value = false;
    }
  }

  async function clearRoster(roomId: string): Promise<void> {
    await apiClearRoster(roomId);
    roster.value = [];
    lastImport.value = null;
  }

  function reset() {
    roster.value = [];
    lastImport.value = null;
    loading.value = false;
  }

  return {
    roster,
    lastImport,
    loading,
    importRoster,
    loadRoster,
    clearRoster,
    reset,
  };
});
