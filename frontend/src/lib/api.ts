import axios from 'axios';
import type {
  AuthResponse,
  GameUnit,
  RosterEntry,
  ImportRosterResponse,
} from '@/types';

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

// Unit placement and management endpoints
export const apiPlaceUnit = (
  roomId: string,
  data: {
    datasheet_id: string;
    model_name: string;
    faction_id: string;
    x: number;
    y: number;
    model_count: number;
    facing_degrees?: number;
    name_on_board?: string;
    owner_player_id?: string;
  }
) => client.post<GameUnit>(`/api/rooms/${roomId}/units`, data);

export const apiMoveUnit = (
  roomId: string,
  unitId: string,
  data: {
    x?: number;
    y?: number;
    facing_degrees?: number;
  }
) =>
  client.patch<GameUnit>(
    `/api/rooms/${roomId}/units/${unitId}`,
    data
  );

export const apiRotateUnit = (
  roomId: string,
  unitId: string,
  facing: number
) =>
  client.patch<GameUnit>(
    `/api/rooms/${roomId}/units/${unitId}`,
    { facing_degrees: facing }
  );

export const apiWoundUnit = (
  roomId: string,
  unitId: string,
  amount: number
) =>
  client.post<GameUnit>(
    `/api/rooms/${roomId}/units/${unitId}/wounds`,
    { amount }
  );

export const apiUpdateUnitStatus = (
  roomId: string,
  unitId: string,
  status: 'alive' | 'in_reserves' | 'dead'
) =>
  client.post<GameUnit>(
    `/api/rooms/${roomId}/units/${unitId}/status`,
    { status }
  );

export const apiRemoveUnit = (roomId: string, unitId: string) =>
  client.delete(`/api/rooms/${roomId}/units/${unitId}`);

export const apiGetRoomUnits = (roomId: string) =>
  client.get<GameUnit[]>(`/api/rooms/${roomId}/units`);

// Wahapedia datasheet endpoints
export interface WhDatasheet {
  id: string;
  name: string;
  faction_id: string;
  role: string;
  [key: string]: any;
}

export interface WhDatasheetModel {
  datasheet_id: string;
  name: string;
  m: string;
  t: string;
  sv: string;
  inv_sv: string;
  w: string;
  ld: string;
  oc: string;
  base_size: string;
  base_size_descr: string;
}

export const apiGetDatasheets = () =>
  client.get<WhDatasheet[]>('/api/wahapedia/datasheets');

export const apiGetDatasheetModels = (datasheetId: string) =>
  client.get<WhDatasheetModel[]>(
    `/api/wahapedia/datasheets/${datasheetId}/models`
  );

// Roster endpoints
export const apiImportRoster = (roomId: string, listforgeJson: unknown) =>
  client.post<ImportRosterResponse>(
    `/api/rooms/${roomId}/roster/import`,
    listforgeJson
  );

export const apiGetRoster = (roomId: string) =>
  client.get<RosterEntry[]>(`/api/rooms/${roomId}/roster`);

export const apiClearRoster = (roomId: string) =>
  client.delete(`/api/rooms/${roomId}/roster`);
