import axios from 'axios';
import type { AuthResponse } from '@/types';

const BASE = import.meta.env.VITE_API_BASE_URL || '';

const client = axios.create({ baseURL: BASE });

// Inject Authorization header on every request
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
});

export const apiRegister = (data: {
  username: string;
  nickname: string;
  password: string;
}) => client.post<AuthResponse>('/api/auth/register', data);

export const apiLogin = (data: { username: string; password: string }) =>
  client.post<AuthResponse>('/api/auth/login', data);

export const apiGetPlayerGames = (id: string) =>
  client.get(`/api/players/${id}/games`);

export const apiCreateRoom = (name: string) =>
  client.post<{ id: string }>('/api/rooms', { name });

export const apiGetRoom = (id: string) => client.get(`/api/rooms/${id}`);

export const apiJoinRoom = (roomId: string, role: 'attacker' | 'defender') =>
  client.post(`/api/rooms/${roomId}/join`, { role });

export const apiStartGame = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/start`);

export const apiPhaseNext = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/phase/next`);

export const apiPhasePrev = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/phase/prev`);

export const apiCloseRoom = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/close`);

export const apiSyncWahapedia = () =>
  client.post('/api/wahapedia/sync');
