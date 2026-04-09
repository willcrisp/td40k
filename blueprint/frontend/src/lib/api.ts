import axios from "axios";
import type { AuthResponse, CounterState } from "@/types";

const BASE = import.meta.env.VITE_API_BASE_URL ?? "";

const client = axios.create({ baseURL: BASE });

client.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const apiRegister = (username: string, password: string) =>
  client.post<AuthResponse>("/api/auth/register", { username, password });

export const apiLogin = (username: string, password: string) =>
  client.post<AuthResponse>("/api/auth/login", { username, password });

export const apiGetCounter = () =>
  client.get<CounterState>("/api/counter");

export const apiIncrementCounter = () =>
  client.post<CounterState>("/api/counter/increment");
