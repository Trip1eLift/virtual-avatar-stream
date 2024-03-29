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

	// 1. Retrieve an unique room_id
	room_id, err := DBW.fetch_unique_room_id()
	if err != nil {
		return err
	}

	// 2. Register room_id with self IP
	ip, err := IP.getIp()
	if err != nil {
		return err
	}
	err = DBW.save_room_id_with_ip(room_id, ip)
	if err != nil {
		return err
	}

	// 3. Feed the room_id to owner client
	err = TM.supply(conn, "Room-Id", room_id)
	if err != nil {
		return err
	}
	log.Printf("Owner host room_id: %s", room_id)

	// 4. Save owner conn for (aisle/guest)
	err = ConnectionCache.addRoom(room_id, conn)
	if err != nil {
		return err
	}

	// 5. Remove room_id from cache when connection closes
	handleClose := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		ConnectionCache.removeRoom(room_id)
		DBW.remove_room_id(room_id)
		return handleClose(code, text)
	})

	// 6. Enter owner reader
	// Dynamically find (aisle/guest) conn
	// - Read from owner and write to (aisle/guest)
	// - Write to owner will be triggered by (aisle/guest)
	proxy_owner_target(room_id)

	return nil
}

func proxy_owner_target(room_id string) error {
	for {
		// wait when there's no target
		log.Println("Owner waiting for target...")
		term, err := ConnectionCache.waitRoom(room_id)
		if err != nil {
			return err
		} else if term {
			// Owner left before target joins
			return nil
		}

		ownerConn, targetConn, fatal, err := ConnectionCache.getRoom(room_id)
		if err != nil {
			if fatal == false {
				continue // If target join and leave
			} else {
				return err // If get room crashes with a fatal error
			}
		}

		for {
			messageType, body, err := ownerConn.ReadMessage()
			if err != nil {
				log.Println(err.Error())
				return err
			}

			// break if target close
			err = targetConn.WriteMessage(messageType, body)
			if err != nil {
				// Handle if target join leave join
				tgRetry, tgErr := ConnectionCache.getTarget(room_id)
				if tgErr == nil && tgRetry.WriteMessage(messageType, body) == nil {
					targetConn = tgRetry
					continue
				} else {
					// Emit first occur error
					log.Println(err.Error())
					break
				}
			}
		}
	}
}
