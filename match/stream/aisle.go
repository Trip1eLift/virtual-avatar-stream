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
	err = ConnectionCache.addTarget(room_id_str, conn)
	if err != nil {
		return err
	}
	go Proxy_target_owner(room_id_str)

	return nil
}

func Proxy_target_owner(room_id string) error {
	ownerConn, targetConn, err := ConnectionCache.getRoom(room_id)
	if err != nil {
		ConnectionCache.removeTarget(room_id)
		return err
	}

	for {
		// TODO: return nil if target close gracefully
		messageType, body, err := targetConn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			ConnectionCache.removeTarget(room_id)
			return err
		}

		err = ownerConn.WriteMessage(messageType, body)
		if err != nil {
			log.Println(err.Error())
			ConnectionCache.removeTarget(room_id)
			return err
		}
	}
}
