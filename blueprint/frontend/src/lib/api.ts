import axios from "axios";
import type { AuthResponse, CounterState, Note } from "@/types";

const BASE = import.meta.env.VITE_API_BASE_URL ?? "";

const client = axios.create({ baseURL: BASE });

// Attach JWT to every request
client.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// On 401, clear session and redirect to login
client.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    const status = (error as { response?: { status?: number } })?.response
      ?.status;
    if (status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user_id");
      localStorage.removeItem("username");
      localStorage.removeItem("is_admin");
      window.location.href = "/auth";
    }
    return Promise.reject(error);
  }
);

export const apiRegister = (username: string, password: string) =>
  client.post<AuthResponse>("/api/auth/register", { username, password });

export const apiLogin = (username: string, password: string) =>
  client.post<AuthResponse>("/api/auth/login", { username, password });

export const apiGetCounter = () => client.get<CounterState>("/api/counter");

export const apiIncrementCounter = () =>
  client.post<CounterState>("/api/counter/increment");

export const apiListNotes = () => client.get<Note[]>("/api/notes");

export const apiCreateNote = (content: string) =>
  client.post<Note>("/api/notes", { content });

export const apiDeleteNote = (id: string) =>
  client.delete<void>(`/api/notes/${id}`);
