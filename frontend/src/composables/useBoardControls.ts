import { ref, type Ref } from 'vue';
import { useBoardStore } from '@/stores/useBoardStore';

const MIN_ZOOM = 0.2;
const MAX_ZOOM = 8.0;

export function useBoardControls(canvas: Ref<HTMLCanvasElement | null>) {
  const board = useBoardStore();
  const isDragging = ref(false);
  let lastX = 0;
  let lastY = 0;

  function onWheel(e: WheelEvent) {
    e.preventDefault();
    if (!canvas.value) return;

    const factor = e.deltaY < 0 ? 1.1 : 0.9;
    const newZoom = Math.min(
      MAX_ZOOM,
      Math.max(MIN_ZOOM, board.zoom * factor)
    );

    const rect = canvas.value.getBoundingClientRect();
    const mx = e.clientX - rect.left;
    const my = e.clientY - rect.top;

    // Zoom towards cursor
    board.panX = mx - ((mx - board.panX) * newZoom) / board.zoom;
    board.panY = my - ((my - board.panY) * newZoom) / board.zoom;
    board.zoom = newZoom;
  }

  function onPointerDown(e: PointerEvent) {
    isDragging.value = true;
    lastX = e.clientX;
    lastY = e.clientY;
    canvas.value?.setPointerCapture(e.pointerId);
  }

  function onPointerMove(e: PointerEvent) {
    if (!isDragging.value) return;
    board.panX += e.clientX - lastX;
    board.panY += e.clientY - lastY;
    lastX = e.clientX;
    lastY = e.clientY;
  }

  function onPointerUp(e: PointerEvent) {
    isDragging.value = false;
    canvas.value?.releasePointerCapture(e.pointerId);
  }

  return { onWheel, onPointerDown, onPointerMove, onPointerUp };
}
