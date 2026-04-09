import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { AuthResponse } from "@/types";

export const usePlayerStore = defineStore("player", () => {
  const token = ref<string | null>(null);
  const playerId = ref<string | null>(null);
  const username = ref<string | null>(null);

  const isAuthenticated = computed(() => !!token.value);

  function init() {
    token.value = localStorage.getItem("token");
    playerId.value = localStorage.getItem("player_id");
    username.value = localStorage.getItem("username");
  }

  function setAuth(data: AuthResponse) {
    token.value = data.token;
    playerId.value = data.player_id;
    username.value = data.username;
    localStorage.setItem("token", data.token);
    localStorage.setItem("player_id", data.player_id);
    localStorage.setItem("username", data.username);
  }

  function logout() {
    token.value = null;
    playerId.value = null;
    username.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("player_id");
    localStorage.removeItem("username");
  }

  return { token, playerId, username, isAuthenticated, init, setAuth, logout };
});
