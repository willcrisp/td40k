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
    header="Create New Game"
    modal
    :style="{ width: '400px' }"
  >
    <div class="flex flex-col gap-4 pt-2">
      <label for="game-name" class="font-semibold">Game Name</label>
      <InputText
        id="game-name"
        v-model="gameName"
        placeholder="e.g. Battle for Cadia"
        @keyup.enter="submit"
        autofocus
      />
    </div>
    <template #footer>
      <Button
        label="Cancel"
        severity="secondary"
        @click="emit('update:visible', false)"
      />
      <Button
        label="Create"
        :disabled="!gameName.trim()"
        @click="submit"
      />
    </template>
  </Dialog>
</template>
