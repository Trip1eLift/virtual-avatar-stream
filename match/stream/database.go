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
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
	psqlconn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
)

func Save_room_id_with_ip(room_id string, ip string) error {
	// room_id_int, err := strconv.Atoi(room_id)
	// if err != nil {
	// 	err = errors.New("Room_id type cast to int error: " + err.Error())
	// 	log.Println(err.Error())
	// 	return err
	// }

	return nil
}

func Remove_room_id(room_id string) error {
	return nil
}

func Fetch_room_id_from_ip(room_id string) (string, error) {
	// room_id_int, err := strconv.Atoi(room_id)
	// if err != nil {
	// 	err = errors.New("Room_id type cast to int error: " + err.Error())
	// 	log.Println(err.Error())
	// 	return "", err
	// }

	return "ip", nil
}

func Fetch_unique_room_id() (string, error) {
	log.Println("postgres:", psqlconn)
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

func Health_database() (string, error) {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		err = errors.New("Postgres connection error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
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
		err = rows.Scan(&room_id, &ip)
		if err != nil {
			err = errors.New("Postgres row scan: " + err.Error())
			log.Println(err.Error())
			return "", err
		}
		results = append(results, []string{room_id, ip})
	}
	res, _ := json.Marshal(results)
	log.Println("Database:\n", string(res))

	return string(res), nil
}
