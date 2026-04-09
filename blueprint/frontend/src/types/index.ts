export interface AuthResponse {
  token: string;
  player_id: string;
  username: string;
  is_admin: boolean;
}

export interface CounterState {
  value: number;
}

export interface Note {
  id: string;
  player_id: string;
  username: string;
  content: string;
  created_at: string;
}

export interface NoteEvent {
  op: "insert" | "delete";
  id: string;
  player_id?: string;
  username?: string;
  content?: string;
  created_at?: string;
}

export interface WsMessage {
  event: string;
  payload: unknown;
}

export interface CounterUpdatePayload {
  value: number;
}
