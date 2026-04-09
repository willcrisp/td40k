CREATE OR REPLACE FUNCTION notify_counter_update()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify(
        'counter_updates',
        json_build_object('value', NEW.value)::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS counter_changed ON counter;
CREATE TRIGGER counter_changed
AFTER UPDATE ON counter
FOR EACH ROW EXECUTE FUNCTION notify_counter_update();
