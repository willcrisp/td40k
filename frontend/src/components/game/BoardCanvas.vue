<script setup lang="ts">
import {
  ref,
  onMounted,
  onUnmounted,
  watchEffect,
  nextTick,
  computed,
} from 'vue';
import { useBoardStore } from '@/stores/useBoardStore';
import { useUnitStore } from '@/stores/useUnitStore';
import { useRoomStore } from '@/stores/useRoomStore';
import { useBoardControls } from '@/composables/useBoardControls';
import {
  apiPlaceUnit,
  apiMoveUnit,
  apiGetRoomUnits,
} from '@/lib/api';
import type { GameUnit } from '@/types';

const BOARD_W_IN = 44;
const BOARD_H_IN = 60;
const BASE_PX = 14; // pixels per inch at zoom 1.0

const canvas = ref<HTMLCanvasElement | null>(null);
const board = useBoardStore();
const unitStore = useUnitStore();
const roomStore = useRoomStore();
const controls = useBoardControls(canvas);

// Dragging state for moving units
const dragging = ref<{
  unitId: string;
  startX: number;
  startY: number;
} | null>(null);

let animFrame: number;

function resizeCanvas() {
  if (!canvas.value) return;
  canvas.value.width = canvas.value.offsetWidth;
  canvas.value.height = canvas.value.offsetHeight;
}

function draw() {
  const el = canvas.value;
  if (!el) return;
  const ctx = el.getContext('2d');
  if (!ctx) return;

  ctx.clearRect(0, 0, el.width, el.height);

  ctx.save();
  ctx.translate(board.panX, board.panY);
  ctx.scale(board.zoom, board.zoom);

  const w = BOARD_W_IN * BASE_PX;
  const h = BOARD_H_IN * BASE_PX;

  // Rotate around the center of the board
  ctx.translate(w / 2, h / 2);
  ctx.rotate((board.rotation * Math.PI) / 180);
  ctx.translate(-w / 2, -h / 2);

  // Board background — dark tactical surface
  ctx.fillStyle = '#131313';
  ctx.fillRect(0, 0, w, h);

  // Grid lines — subtle industrial grid
  ctx.strokeStyle = 'rgba(229, 226, 225, 0.06)';
  ctx.lineWidth = 0.5 / board.zoom;

  for (let x = 0; x <= BOARD_W_IN; x++) {
    ctx.beginPath();
    ctx.moveTo(x * BASE_PX, 0);
    ctx.lineTo(x * BASE_PX, h);
    ctx.stroke();
  }

  for (let y = 0; y <= BOARD_H_IN; y++) {
    ctx.beginPath();
    ctx.moveTo(0, y * BASE_PX);
    ctx.lineTo(w, y * BASE_PX);
    ctx.stroke();
  }

  // Board border — muted industrial edge
  ctx.strokeStyle = '#5a403d';
  ctx.lineWidth = 2 / board.zoom;
  ctx.strokeRect(0, 0, w, h);

  // Inch labels at 6" intervals — amber phosphor style
  ctx.fillStyle = 'rgba(255, 185, 82, 0.25)';
  ctx.font = `${10 / board.zoom}px 'Space Grotesk', monospace`;
  
  // Width labels (top edge)
  for (let x = 0; x <= BOARD_W_IN; x += 6) {
    ctx.fillText(`${x}"`, x * BASE_PX + 2, 10 / board.zoom);
  }

  // Height labels (left edge)
  for (let y = 6; y <= BOARD_H_IN; y += 6) {
    ctx.fillText(`${y}"`, 2, y * BASE_PX - 2);
  }

  // Draw units
  unitStore.units.forEach((unit) => {
    drawUnit(ctx, unit, BASE_PX, unitStore.selectedUnitId === unit.id);
  });

  ctx.restore();
}

function drawUnit(
  ctx: CanvasRenderingContext2D,
  unit: GameUnit,
  basePx: number,
  isSelected: boolean
) {
  const x = unit.x * basePx;
  const y = unit.y * basePx;

  // Determine color based on owner
  let color = 'rgba(100, 100, 100, 0.5)'; // default gray
  // TODO: Color based on attacker/defender from room state

  // Determine shape and draw footprint (base size in mm)
  const baseSize = parseFloat(unit.model_name) || 32; // fallback
  const radiusInches = baseSize / 25.4 / 2; // convert mm to inches to pixels
  const radiusPx = radiusInches * basePx;

  ctx.save();
  ctx.translate(x, y);
  ctx.rotate((unit.facing_degrees * Math.PI) / 180);

  // Draw base shape (assume circle for now)
  ctx.fillStyle = color;
  ctx.beginPath();
  ctx.arc(0, 0, radiusPx, 0, 2 * Math.PI);
  ctx.fill();

  // Draw border
  ctx.strokeStyle = isSelected ? 'rgba(255, 255, 0, 1)' : color;
  ctx.lineWidth = isSelected ? 2 : 1;
  ctx.stroke();

  // Draw facing indicator (line at top)
  ctx.strokeStyle = color;
  ctx.lineWidth = 2;
  ctx.beginPath();
  ctx.moveTo(0, 0);
  ctx.lineTo(0, -radiusPx * 0.7);
  ctx.stroke();

  // Draw label
  ctx.fillStyle = 'rgba(255, 255, 255, 1)';
  ctx.font = `bold ${12 / board.zoom}px Arial`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(`${unit.model_count}`, 0, 0);

  ctx.restore();
}

// Convert canvas pixel coordinates to board inches
function canvasToBoard(
  canvasX: number,
  canvasY: number
): { x: number; y: number } {
  const w = BOARD_W_IN * BASE_PX;
  const h = BOARD_H_IN * BASE_PX;

  // Apply inverse transformations
  let x = (canvasX - board.panX) / board.zoom;
  let y = (canvasY - board.panY) / board.zoom;

  // Undo rotation
  const angle = (-board.rotation * Math.PI) / 180;
  const cx = w / 2;
  const cy = h / 2;

  x -= cx;
  y -= cy;

  const cosA = Math.cos(angle);
  const sinA = Math.sin(angle);
  const rotX = x * cosA - y * sinA;
  const rotY = x * sinA + y * cosA;

  x = rotX + cx;
  y = rotY + cy;

  // Convert to inches
  return {
    x: x / BASE_PX,
    y: y / BASE_PX,
  };
}

// Check if point is inside a unit's base
function hitTestUnit(
  boardX: number,
  boardY: number,
  unit: GameUnit
): boolean {
  const dx = boardX - unit.x;
  const dy = boardY - unit.y;
  // Simple circle hit test; in reality should use base_size
  const radius = 1.0; // 1 inch radius approximate
  return Math.sqrt(dx * dx + dy * dy) <= radius;
}

// Handle canvas click for unit placement or selection
async function handleCanvasClick(
  e: MouseEvent
): Promise<void> {
  if (!canvas.value) return;

  const rect = canvas.value.getBoundingClientRect();
  const canvasX = e.clientX - rect.left;
  const canvasY = e.clientY - rect.top;
  const boardCoords = canvasToBoard(canvasX, canvasY);

  // Check if we're placing a new unit
  if (unitStore.placingUnitType) {
    if (!roomStore.roomId) return;

    try {
      const placed = await apiPlaceUnit(roomStore.roomId, {
        datasheet_id: unitStore.placingUnitType.datasheetId,
        model_name: unitStore.placingUnitType.modelName,
        faction_id: unitStore.placingUnitType.factionId,
        x: boardCoords.x,
        y: boardCoords.y,
        model_count: unitStore.placingUnitType.modelCount,
        facing_degrees: 0,
      });

      unitStore.addUnit(placed.data);
      unitStore.clearPlacing();
      return;
    } catch (err) {
      console.error('Failed to place unit:', err);
      return;
    }
  }

  // Check if clicking on an existing unit
  for (const unit of unitStore.units) {
    if (hitTestUnit(boardCoords.x, boardCoords.y, unit)) {
      unitStore.selectUnit(unit.id);
      return;
    }
  }

  // Clicked on empty space
  unitStore.selectUnit(null);
}

function loop() {
  draw();
  animFrame = requestAnimationFrame(loop);
}

onMounted(async () => {
  await nextTick();
  resizeCanvas();
  window.addEventListener('resize', resizeCanvas);
  // Centre the board
  if (canvas.value) {
    board.panX = (canvas.value.width - BOARD_W_IN * BASE_PX) / 2;
    board.panY = (canvas.value.height - BOARD_H_IN * BASE_PX) / 2;
  }
  // Load units
  if (roomStore.roomId) {
    try {
      const { data } = await apiGetRoomUnits(roomStore.roomId);
      unitStore.setUnits(data);
    } catch (err) {
      console.error('Failed to load units:', err);
    }
  }
  loop();
});

onUnmounted(() => {
  cancelAnimationFrame(animFrame);
  window.removeEventListener('resize', resizeCanvas);
});

// Redraw whenever store changes (zoom/pan updates trigger a new frame anyway
// via rAF loop, but watchEffect ensures reactivity is tracked)
watchEffect(() => {
  void board.zoom;
  void board.panX;
  void board.panY;
});
</script>

<template>
  <canvas
    ref="canvas"
    class="w-full h-full block cursor-grab active:cursor-grabbing"
    @wheel="controls.onWheel"
    @pointerdown="controls.onPointerDown"
    @pointermove="controls.onPointerMove"
    @pointerup="controls.onPointerUp"
    @click="handleCanvasClick"
  />
</template>
