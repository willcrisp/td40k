import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { GameUnit } from '@/types';
import { usePlayerStore } from './usePlayerStore';

export const useUnitStore = defineStore('units', () => {
  const units = ref<GameUnit[]>([]);
  const selectedUnitId = ref<string | null>(null);
  const placingUnitType = ref<{
    datasheetId: string;
    modelName: string;
    modelCount: number;
    factionId: string;
  } | null>(null);

  const selectedUnit = computed(() => {
    if (!selectedUnitId.value) return null;
    return units.value.find((u) => u.id === selectedUnitId.value);
  });

  const playerStore = usePlayerStore();

  const playerUnits = computed(() => {
    if (!playerStore.playerId) return [];
    return units.value.filter((u) => u.owner_player_id === playerStore.playerId);
  });

  const alivePlayers = computed(() =>
    units.value.filter((u) => u.status === 'alive')
  );

  const reserveUnits = computed(() =>
    units.value.filter((u) => u.status === 'in_reserves')
  );

  const deadUnits = computed(() =>
    units.value.filter((u) => u.status === 'dead')
  );

  // Initialize units from server
  function setUnits(newUnits: GameUnit[]) {
    units.value = newUnits;
  }

  // Add single unit
  function addUnit(unit: GameUnit) {
    const existing = units.value.findIndex((u) => u.id === unit.id);
    if (existing >= 0) {
      units.value[existing] = unit;
    } else {
      units.value.push(unit);
    }
  }

  // Update unit properties
  function updateUnit(unitId: string, updates: Partial<GameUnit>) {
    const idx = units.value.findIndex((u) => u.id === unitId);
    if (idx >= 0) {
      units.value[idx] = { ...units.value[idx], ...updates };
    }
  }

  // Remove unit
  function removeUnit(unitId: string) {
    units.value = units.value.filter((u) => u.id !== unitId);
    if (selectedUnitId.value === unitId) {
      selectedUnitId.value = null;
    }
  }

  // Select unit for inspection
  function selectUnit(unitId: string | null) {
    selectedUnitId.value = unitId;
  }

  // Start placing a new unit
  function startPlacingUnit(
    datasheetId: string,
    modelName: string,
    modelCount: number,
    factionId: string
  ) {
    placingUnitType.value = {
      datasheetId,
      modelName,
      modelCount,
      factionId,
    };
  }

  // Cancel unit placement
  function cancelPlacing() {
    placingUnitType.value = null;
  }

  // Clear placing state (after successful placement)
  function clearPlacing() {
    placingUnitType.value = null;
  }

  return {
    units,
    selectedUnitId,
    selectedUnit,
    placingUnitType,
    playerUnits,
    alivePlayers,
    reserveUnits,
    deadUnits,
    setUnits,
    addUnit,
    updateUnit,
    removeUnit,
    selectUnit,
    startPlacingUnit,
    cancelPlacing,
    clearPlacing,
  };
});
