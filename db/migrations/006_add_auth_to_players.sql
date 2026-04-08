-- Migration 006: Add username/password authentication to players.
-- Old sessions used client-generated UUIDs (X-Player-ID header) which are
-- now invalid. Truncate dependent tables so the schema change is clean.
TRUNCATE room_events, rooms, players;

ALTER TABLE players
  ADD COLUMN username      TEXT UNIQUE NOT NULL,
  ADD COLUMN password_hash TEXT        NOT NULL;
