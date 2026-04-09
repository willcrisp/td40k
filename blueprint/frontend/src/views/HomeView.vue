<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import { useCounterStore } from "@/stores/useCounterStore";
import { useWebSocketStore } from "@/stores/useWebSocketStore";
import { usePlayerStore } from "@/stores/usePlayerStore";

const router = useRouter();
const counterStore = useCounterStore();
const wsStore = useWebSocketStore();
const playerStore = usePlayerStore();

const { value, loading } = storeToRefs(counterStore);
const { username } = storeToRefs(playerStore);

onMounted(async () => {
  await counterStore.fetchCounter();
  wsStore.connect();
});

onUnmounted(() => {
  wsStore.disconnect();
});

function handleLogout() {
  wsStore.disconnect();
  playerStore.logout();
  router.push("/auth");
}
</script>

<template>
  <div class="flex flex-column align-items-center gap-4 p-6">
    <Toolbar style="width: 100%; max-width: 480px">
      <template #start>
        <span class="font-semibold">{{ username }}</span>
      </template>
      <template #end>
        <Button label="Logout" severity="secondary" size="small" @click="handleLogout" />
      </template>
    </Toolbar>

    <Card style="width: 100%; max-width: 480px">
      <template #title>Shared Counter</template>
      <template #content>
        <div class="flex flex-column align-items-center gap-4">
          <div class="text-8xl font-bold">{{ value }}</div>
          <Button
            label="Increment"
            :loading="loading"
            size="large"
            @click="counterStore.increment()"
          />
        </div>
      </template>
    </Card>
  </div>
</template>
