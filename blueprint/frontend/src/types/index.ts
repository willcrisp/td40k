export interface AuthResponse {
  token: string;
  player_id: string;
  username: string;
}

export interface CounterState {
  value: number;
}

export interface WsMessage {
  event: string;
  payload: unknown;
}

export interface CounterUpdatePayload {
  value: number;
}
