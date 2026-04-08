export type Phase =
  | 'command'
  | 'movement'
  | 'shooting'
  | 'charge'
  | 'fight';

export type RoomStatus = 'lobby' | 'active' | 'finished' | 'closed';

export type PlayerRole =
  | 'attacker'
  | 'defender'
  | 'game_master'
  | null;

export type ActivePlayer = 'attacker' | 'defender';

export interface Room {
  id: string;
  name: string;
  status: RoomStatus;
  game_master_id: string;
  attacker_id: string | null;
  defender_id: string | null;
  battle_round: number;
  active_player: ActivePlayer;
  current_phase: Phase;
  winner: ActivePlayer | null;
  created_at: string;
  updated_at: string;
}

export interface OwnedGameSummary {
  id: string;
  name: string;
  status: RoomStatus;
  battle_round: number;
  active_player: ActivePlayer;
  current_phase: Phase;
  attacker_id: string | null;
  defender_id: string | null;
  created_at: string;
}

export interface JoinedGameSummary {
  id: string;
  name: string;
  status: RoomStatus;
  role: 'attacker' | 'defender';
  battle_round: number;
  current_phase: Phase;
  created_at: string;
}

export interface RoomStatePayload {
  room_id: string;
  name: string;
  status: RoomStatus;
  battle_round: number;
  active_player: ActivePlayer;
  current_phase: Phase;
  winner: ActivePlayer | null;
  attacker_id: string | null;
  defender_id: string | null;
  game_master_id: string;
}

export interface WsMessage {
  event: 'room_state';
  payload: RoomStatePayload;
}

export const PHASES: Phase[] = [
  'command',
  'movement',
  'shooting',
  'charge',
  'fight',
];

export const PHASE_LABELS: Record<Phase, string> = {
  command: 'Command Phase',
  movement: 'Movement Phase',
  shooting: 'Shooting Phase',
  charge: 'Charge Phase',
  fight: 'Fight Phase',
};

export const PHASE_NUMBERS: Record<Phase, number> = {
  command: 1,
  movement: 2,
  shooting: 3,
  charge: 4,
  fight: 5,
};

export interface AuthResponse {
  token: string;
  player_id: string;
  username: string;
  nickname: string;
  is_admin: boolean;
}
