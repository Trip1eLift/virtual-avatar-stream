package stream

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func HandleAisle(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "aisle" {
		return nil
	}

	aisle_key, err := Demand(conn, "Authorization")
	if err != nil {
		return err
	}
	if aisle_key != os.Getenv("AISLE_KEY") {
		err = errors.New("[Critical] Aisle key is incorrect.")
		log.Println(err)
		return err
	}
	room_id_str, err := Demand(conn, "Room-Id")
	if err != nil {
		return err
	}
	log.Printf("Aisle target room_id: %s\n", room_id_str)

	// Find owner conn
	// Enter aisle-owner reader
	// - Read from aisle and write to owner
	// - Write to aisle will be triggered by owner

	return nil
}
