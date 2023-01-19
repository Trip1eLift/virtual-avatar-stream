package stream

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbhost     = os.Getenv("DB_HOST")
	dbport     = os.Getenv("DB_PORT")
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASS")
	dbname     = os.Getenv("DB_NAME")
	psqlconn   = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)
)

type Database struct {
}

// TODO PRIO: migrate to pgx
//            docs: https://pkg.go.dev/github.com/jackc/pgx/v5
//            repo: https://github.com/jackc/pgx
// It supports multi-statements: https://stackoverflow.com/questions/38998267/how-to-execute-a-sql-file

func (d *Database) save_room_id_with_ip(room_id string, ip string) error {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO rooms(room_id, task_private_ip) VALUES($1, $2);", room_id, ip)
	if err != nil {
		err = errors.New("Postgres exec error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	return nil
}

func (d *Database) remove_room_id(room_id string) error {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM rooms WHERE room_id=$1;", room_id)
	if err != nil {
		err = errors.New("Postgres exec error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	return nil
}

func (d *Database) fetch_ip_from_room_id(room_id string) (string, bool, error) {
	// bool marks if the error is fatal
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return "", true, err
	}
	defer db.Close()

	var ip string
	err = db.QueryRow("SELECT task_private_ip FROM rooms WHERE room_id=$1;", room_id).Scan(&ip)
	switch err {
	case sql.ErrNoRows:
		err = errors.New(fmt.Sprintf("IP of room_id: %s not found. %s", room_id, err.Error()))
		log.Println(err.Error())
		return "", false, err
	case nil:
		return ip, false, nil
	default:
		err = errors.New("Postgres query unexpected error: " + err.Error())
		log.Println(err.Error())
		return "", true, err
	}
}

func (d *Database) fetch_unique_room_id() (string, error) {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	defer db.Close()

	var room_id string
	err = db.QueryRow("SELECT nextval('room_id_seq');").Scan(&room_id)
	if err != nil {
		err = errors.New("Postgres query error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	log.Println("room_id:", room_id)
	return room_id, nil
}

func (d *Database) health_database() (string, error) {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM rooms;")
	if err != nil {
		err = errors.New("Postgres query error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	defer rows.Close()

	results := make([][]string, 0)
	for rows.Next() {
		var room_id string
		var ip string
		var timestamp string
		err = rows.Scan(&room_id, &ip, &timestamp)
		if err != nil {
			err = errors.New("Postgres row scan: " + err.Error())
			log.Println(err.Error())
			return "", err
		}
		results = append(results, []string{room_id, ip, timestamp})
	}
	res, _ := json.Marshal(results)
	log.Println("Database:\n", string(res))

	return string(res), nil
}

func (d *Database) fetch_an_non_self_ip(self_ip string) (string, error) {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	defer db.Close()

	var target_ip string
	err = db.QueryRow("SELECT task_private_ip FROM rooms WHERE task_private_ip!=$1 LIMIT 1;", self_ip).Scan(&target_ip)
	if err != nil {
		err = errors.New("Postgres query error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}

	return target_ip, nil
}

// TODO: add a database init sequence
func (d *Database) initialize(self_ip string) error {
	// TODO: try to execute a sql file here
	return nil
}

var DB = Database{}
