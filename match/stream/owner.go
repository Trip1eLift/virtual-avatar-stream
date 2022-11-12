package stream

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleOwner(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "owner" {
		return nil
	}

	// 1. Retrieve a unique room_id
	// TODO: check database to find unique room_id
	// TODO: register room_id with IP to the database
	room_id := "1"

	// 2. Feed the room_id to owner client
	err := Supply(conn, "Room-Id", room_id)
	if err != nil {
		return err
	}

	// 3. Save owner conn for (aisle/guest)
	err = ConnectionCache.addRoom(room_id, conn)
	if err != nil {
		return err
	}

	// 4. Enter owner reader
	// Dynamically find (aisle/guest) conn
	// - Read from owner and write to (aisle/guest)
	// - Write to owner will be triggered by (aisle/guest)
	err = reader(room_id)
	if err != nil {
		return err
	}

	// 5. Remove room_id from cache when connection closes
	// TODO: how do detect connection closure

	return nil
}

func reader(room_id string) error {
	if _, found := ConnectionCache.index[room_id]; found == false {
		err := errors.New(fmt.Sprintf("Owner reader failure. Room_id: %s does not exist.", room_id))
		log.Println(err.Error())
		return err
	}

	for {
		// TODO: handle client close here?
		messageType, body, err := ConnectionCache.index[room_id].owner.ReadMessage()
		if err != nil {
			log.Println(err)
			return err
		}

		if ConnectionCache.index[room_id].target != nil {
			err = ConnectionCache.index[room_id].target.WriteMessage(messageType, body)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
}
