<script setup lang="ts">
import { onMounted } from "vue";
import { useUserStore } from "@/stores/useUserStore";

const userStore = useUserStore();

onMounted(() => {
  userStore.init();
});
</script>

<template>
  <Toast />
  <RouterView />
</template>
