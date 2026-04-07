<script setup lang="ts">
import type { OwnedGameSummary } from '@/types';
import { PHASE_LABELS } from '@/types';
import Button from 'primevue/button';
import Card from 'primevue/card';
import Tag from 'primevue/tag';

const props = defineProps<{ game: OwnedGameSummary }>();
const emit = defineEmits<{
  (e: 'enter', id: string): void;
  (e: 'close', id: string): void;
}>();

const statusSeverity: Record<string, string> = {
  lobby: 'warn',
  active: 'success',
  finished: 'info',
  closed: 'secondary',
};
</script>

<template>
  <Card>
    <template #title>{{ props.game.name }}</template>
    <template #subtitle>
      <Tag
        :value="props.game.status.toUpperCase()"
        :severity="statusSeverity[props.game.status]"
      />
    </template>
    <template #content>
      <p v-if="props.game.status === 'active'" class="text-sm">
        Round {{ props.game.battle_round }} of 5 &mdash;
        {{ props.game.active_player.toUpperCase() }}&apos;s Turn
        <br />
        {{ PHASE_LABELS[props.game.current_phase] }}
      </p>
      <p v-else-if="props.game.status === 'lobby'" class="text-sm">
        Waiting for players&hellip;
        <br />
        Attacker:
        {{ props.game.attacker_id ? '✅ Joined' : '⏳ Empty' }}
        &nbsp;|&nbsp; Defender:
        {{ props.game.defender_id ? '✅ Joined' : '⏳ Empty' }}
      </p>
      <p v-else-if="props.game.status === 'finished'" class="text-sm">
        Game complete
      </p>
    </template>
    <template #footer>
      <div class="flex gap-2">
        <Button
          label="Enter"
          size="small"
          @click="emit('enter', props.game.id)"
        />
        <Button
          label="Close Game"
          size="small"
          severity="danger"
          outlined
          :disabled="
            props.game.status === 'closed' ||
            props.game.status === 'finished'
          "
          @click="emit('close', props.game.id)"
        />
      </div>
    </template>
  </Card>
</template>
