package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleOwner(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "owner" {
		return nil
	}
	// TODO: check database to find unique room_id
	// TODO: register room_id with IP to the database
	message := map[string]interface{}{
		"key":   "room_id",
		"value": 1,
	}
	if payload, err := json.Marshal(message); err != nil {
		return err
	} else {
		conn.WriteMessage(websocket.TextMessage, payload)

		// save owner conn for (aisle/guest)
		// Enter owner reader
		// Dynamically find (aisle/guest) conn
		// - Read from owner and write to (aisle/guest)
		// - Write to owner will be triggered by (aisle/guest)
	}
	return nil
}

func reader(conn *websocket.Conn) {
	for {
		// TODO: handle client close here?

		messageType, body, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(body))

		if err := conn.WriteMessage(messageType, body); err != nil {
			log.Println(err)
			return
		}
	}
}
