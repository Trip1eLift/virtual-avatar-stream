package stream

import (
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
	// 5. Remove room_id from cache when connection closes
	go proxy_owner_target(room_id)

	return nil
}

func proxy_owner_target(room_id string) error {
	for {
		// wait when there's no target
		log.Println("HERE wait")
		err := ConnectionCache.waitRoom(room_id)
		log.Println("HERE after wait")
		if err != nil {
			ConnectionCache.removeRoom(room_id)
			return err
		}
		ownerConn, targetConn, err := ConnectionCache.getRoom(room_id)
		if err != nil {
			ConnectionCache.removeRoom(room_id)
			return err
		}

		for {

			messageType, body, err := ownerConn.ReadMessage()
			if err != nil {
				log.Println(err.Error())
				ConnectionCache.removeRoom(room_id)
				return err
			}

			// break if target close
			err = targetConn.WriteMessage(messageType, body)
			if err != nil {
				log.Println(err.Error())
				break
			}
		}
	}
}