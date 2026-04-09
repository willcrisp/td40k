import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { AuthResponse } from "@/types";

export const useUserStore = defineStore("user", () => {
  const token = ref<string | null>(null);
  const userId = ref<string | null>(null);
  const username = ref<string | null>(null);
  const isAdmin = ref<boolean>(false);

  const isAuthenticated = computed(() => !!token.value);

  function init() {
    token.value = localStorage.getItem("token");
    userId.value = localStorage.getItem("user_id");
    username.value = localStorage.getItem("username");
    isAdmin.value = localStorage.getItem("is_admin") === "true";
  }

  function setAuth(data: AuthResponse) {
    token.value = data.token;
    userId.value = data.user_id;
    username.value = data.username;
    isAdmin.value = data.is_admin;
    localStorage.setItem("token", data.token);
    localStorage.setItem("user_id", data.user_id);
    localStorage.setItem("username", data.username);
    localStorage.setItem("is_admin", String(data.is_admin));
  }

  function logout() {
    token.value = null;
    userId.value = null;
    username.value = null;
    isAdmin.value = false;
    localStorage.removeItem("token");
    localStorage.removeItem("user_id");
    localStorage.removeItem("username");
    localStorage.removeItem("is_admin");
  }

  return {
    token,
    userId,
    username,
    isAdmin,
    isAuthenticated,
    init,
    setAuth,
    logout,
  };
});
