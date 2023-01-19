CREATE TABLE IF NOT EXISTS rooms(
  room_id VARCHAR(10) PRIMARY KEY, 
  task_private_ip VARCHAR(20) NOT NULL,
  timestamp timestamp NOT NULL DEFAULT NOW()
);

CREATE FUNCTION expire_rooms() RETURNS trigger
  LANGUAGE plpgsql
  AS $$
    BEGIN
      DELETE FROM rooms WHERE timestamp < NOW() - INTERVAL '24 hours';
      RETURN NEW;
    END;
  $$;

CREATE TRIGGER expire_rooms_trigger
  AFTER INSERT ON rooms
  EXECUTE PROCEDURE expire_rooms();

CREATE SEQUENCE IF NOT EXISTS room_id_seq
  START WITH 1
  INCREMENT BY 1
  MINVALUE 1
  MAXVALUE 999999
  CACHE 1;

-- INSERT INTO rooms(room_id, task_private_ip) VALUES('-1', 'test-ip-1');
-- INSERT INTO rooms(room_id, task_private_ip) VALUES('-2', 'test-ip-2');