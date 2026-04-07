<script setup lang="ts">
import type { JoinedGameSummary } from '@/types';
import { PHASE_LABELS } from '@/types';
import Button from 'primevue/button';
import Card from 'primevue/card';
import Tag from 'primevue/tag';

const props = defineProps<{ game: JoinedGameSummary }>();
const emit = defineEmits<{
  (e: 'rejoin', id: string): void;
}>();

const statusSeverity: Record<string, string> = {
  lobby: 'warn',
  active: 'success',
  finished: 'info',
  closed: 'secondary',
};

const roleLabel: Record<string, string> = {
  attacker: 'Attacker',
  defender: 'Defender',
};
</script>

<template>
  <Card>
    <template #title>{{ props.game.name }}</template>
    <template #subtitle>
      <div class="flex gap-2">
        <Tag
          :value="props.game.status.toUpperCase()"
          :severity="statusSeverity[props.game.status]"
        />
        <Tag :value="roleLabel[props.game.role]" severity="info" />
      </div>
    </template>
    <template #content>
      <p class="text-sm">
        Round {{ props.game.battle_round }} of 5
        <br />
        {{ PHASE_LABELS[props.game.current_phase] }}
      </p>
    </template>
    <template #footer>
      <Button
        label="Rejoin"
        size="small"
        :disabled="props.game.status === 'closed'"
        @click="emit('rejoin', props.game.id)"
      />
    </template>
  </Card>
</template>
