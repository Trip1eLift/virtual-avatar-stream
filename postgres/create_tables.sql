CREATE TABLE IF NOT EXISTS rooms(room_id VARCHAR(10) PRIMARY KEY, task_private_ip VARCHAR(20));

CREATE SEQUENCE IF NOT EXISTS room_id_seq
  START WITH 1
  INCREMENT BY 1
  MINVALUE 1
  MAXVALUE 999999
  CACHE 1;

INSERT INTO rooms(room_id, task_private_ip) VALUES('-1', 'test-ip-1');
INSERT INTO rooms(room_id, task_private_ip) VALUES('-2', 'test-ip-2');