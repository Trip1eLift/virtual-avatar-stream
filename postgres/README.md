# Postgres Database Commands
https://www.postgresql.org/docs/current/app-psql.html

## Login database
```shell
psql -U postgres_user
```

## List databases
```shell
\l
```

## List tables
```shell
\d
```

## Query
```shell
SELECT * FROM rooms;

SELECT nextval('room_id_seq');
```