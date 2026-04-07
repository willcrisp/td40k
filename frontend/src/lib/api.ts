import axios from 'axios';

const BASE = import.meta.env.VITE_API_BASE_URL || '';

const client = axios.create({ baseURL: BASE });

// Inject X-Player-ID on every request
client.interceptors.request.use((config) => {
  const playerId = localStorage.getItem('player_id');
  if (playerId) {
    config.headers['X-Player-ID'] = playerId;
  }
  return config;
});

export const apiUpsertPlayer = (id: string, nickname: string) =>
  client.post('/api/players', { id, nickname });

export const apiGetPlayerGames = (id: string) =>
  client.get(`/api/players/${id}/games`);

export const apiCreateRoom = (name: string) =>
  client.post<{ id: string }>('/api/rooms', { name });

export const apiGetRoom = (id: string) =>
  client.get(`/api/rooms/${id}`);

export const apiJoinRoom = (
  roomId: string,
  role: 'attacker' | 'defender'
) => client.post(`/api/rooms/${roomId}/join`, { role });

export const apiStartGame = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/start`);

export const apiPhaseNext = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/phase/next`);

export const apiPhasePrev = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/phase/prev`);

export const apiCloseRoom = (roomId: string) =>
  client.post(`/api/rooms/${roomId}/close`);
