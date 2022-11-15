package stream

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
	psqlconn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
)

func Fetch_room_id_from_ip(room_id string) (string, error) {
	// room_id_int, err := strconv.Atoi(room_id)
	// if err != nil {
	// 	err = errors.New("Room_id type cast to int error: " + err.Error())
	// 	log.Println(err.Error())
	// 	return "", err
	// }

	return "ip", nil
}

func Save_room_id_with_ip(room_id string, ip string) error {
	// room_id_int, err := strconv.Atoi(room_id)
	// if err != nil {
	// 	err = errors.New("Room_id type cast to int error: " + err.Error())
	// 	log.Println(err.Error())
	// 	return err
	// }

	return nil
}

func Fetch_unique_room_id() (string, error) {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	rows, err := db.Query("SELECT nextval('room_id_seq');")
	if err != nil {
		err = errors.New("Postgres query error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	log.Println("rows:", rows)
	return "1", nil
}

func Health_database() error {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	rows, err := db.Query("SELECT * FROM rooms;")
	if err != nil {
		err = errors.New("Postgres query error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	log.Println("Database health state:\n", rows)
	return nil
}
