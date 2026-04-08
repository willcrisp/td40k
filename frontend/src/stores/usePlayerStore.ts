import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { AuthResponse } from '@/types';

export const usePlayerStore = defineStore('player', () => {
  const token = ref<string>('');
  const playerId = ref<string>('');
  const username = ref<string>('');
  const nickname = ref<string>('');
  const isAdmin = ref<boolean>(false);
  const initialized = ref(false);

  const isAuthenticated = computed(() => !!token.value);

  function init() {
    token.value = localStorage.getItem('token') ?? '';
    playerId.value = localStorage.getItem('player_id') ?? '';
    username.value = localStorage.getItem('username') ?? '';
    nickname.value = localStorage.getItem('nickname') ?? '';
    isAdmin.value = localStorage.getItem('is_admin') === 'true';
    initialized.value = true;
  }

  function setAuth(data: AuthResponse) {
    token.value = data.token;
    playerId.value = data.player_id;
    username.value = data.username;
    nickname.value = data.nickname;
    isAdmin.value = data.is_admin;
    localStorage.setItem('token', data.token);
    localStorage.setItem('player_id', data.player_id);
    localStorage.setItem('username', data.username);
    localStorage.setItem('nickname', data.nickname);
    localStorage.setItem('is_admin', String(data.is_admin));
  }

  function logout() {
    token.value = '';
    playerId.value = '';
    username.value = '';
    nickname.value = '';
    isAdmin.value = false;
    localStorage.removeItem('token');
    localStorage.removeItem('player_id');
    localStorage.removeItem('username');
    localStorage.removeItem('nickname');
    localStorage.removeItem('is_admin');
  }

  return {
    token,
    playerId,
    username,
    nickname,
    isAdmin,
    initialized,
    isAuthenticated,
    init,
    setAuth,
    logout,
  };
});
