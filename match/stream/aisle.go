package stream

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// TODO: Fix guest->aisle->owner guest join leave join

func HandleAisle(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "aisle" {
		return nil
	}

	// 1. Establish guest-aisle connection, check authorization, and retrieve room_id
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
	log.Printf("Aisle join room_id: %s\n", room_id_str)

	// 2. Save aisle conn for owner
	err = ConnectionCache.addTarget(room_id_str, conn)
	if err != nil {
		return err
	}

	// 3. Remove aisle from cache when connection closes
	handleClose := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		ConnectionCache.removeTarget(room_id_str)
		return handleClose(code, text)
	})

	// 4. Enter aisle-owner reader
	// - Read from aisle and write to owner
	// - Write to aisle will be triggered by owner
	Proxy_target_owner(room_id_str)

	return nil
}

func Proxy_target_owner(room_id string) error {
	ownerConn, targetConn, _, err := ConnectionCache.getRoom(room_id)
	if err != nil {
		return err
	}

	for {
		// TODO: return nil if target close gracefully
		messageType, body, err := targetConn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return err
		}

		err = ownerConn.WriteMessage(messageType, body)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
}
