-- +goose Up
-- Initialize the counter table
CREATE TABLE IF NOT EXISTS counter (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    value INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert the default counter
INSERT INTO counter (name, value) VALUES ('main', 0) ON CONFLICT DO NOTHING;

-- +goose StatementBegin
-- Create the realtime Postgres notification trigger function
CREATE OR REPLACE FUNCTION notify_counter_update()
RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify(
        'counter_updates',
        json_build_object(
            'name', NEW.name,
            'value', NEW.value,
            'updated_at', NEW.updated_at
        )::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Ensure the trigger attaches to the table
CREATE OR REPLACE TRIGGER counter_update_trigger
AFTER UPDATE ON counter
FOR EACH ROW
EXECUTE FUNCTION notify_counter_update();

-- +goose Down
DROP TRIGGER IF EXISTS counter_update_trigger ON counter;
DROP FUNCTION IF EXISTS notify_counter_update();
DROP TABLE IF EXISTS counter;
