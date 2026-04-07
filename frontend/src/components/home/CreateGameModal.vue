<script setup lang="ts">
import { ref } from 'vue';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';

const props = defineProps<{ visible: boolean }>();
const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void;
  (e: 'create', name: string): void;
}>();

const gameName = ref('');

function submit() {
  if (!gameName.value.trim()) return;
  emit('create', gameName.value.trim());
  gameName.value = '';
}
</script>

<template>
  <Dialog
    :visible="props.visible"
    @update:visible="emit('update:visible', $event)"
    modal
    :style="{ width: '450px' }"
    :closable="false"
  >
    <template #header>
      <div class="flex items-center gap-2">
        <div class="w-1.5 h-6 bg-tertiary"></div>
        <h2 class="text-xl font-display text-white">NEW GAME</h2>
      </div>
    </template>

    <div class="flex flex-col gap-6 py-6 riveted">
      <div class="flex flex-col gap-2">
        <label for="game-name" class="text-xs font-mono text-surface-variant uppercase">
          Game name
        </label>
        <InputText
          id="game-name"
          v-model="gameName"
          placeholder="Enter a name for your game..."
          class="font-mono text-sm"
          @keyup.enter="submit"
          autofocus
        />
        <p class="text-[10px] font-mono text-tertiary mt-1 opacity-70">
          Note: This name will be visible to all players.
        </p>
      </div>
    </div>

    <template #footer>
      <div class="flex gap-4 w-full pt-4">
        <Button
          label="CANCEL"
          class="btn-secondary-tactical flex-1"
          @click="emit('update:visible', false)"
        />
        <Button
          label="CREATE"
          class="btn-tactical flex-1"
          :disabled="!gameName.trim()"
          @click="submit"
        />
      </div>
    </template>
  </Dialog>
</template>
