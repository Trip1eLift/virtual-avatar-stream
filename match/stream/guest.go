package stream

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func HandleGuest(conn *websocket.Conn, request *http.Request, port string) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "guest" {
		return nil
	}

	// 1. Establish client-guest connection and retrieve room_id as int
	room_id, err := Demand(conn, "Room-Id")
	if err != nil {
		return err
	}
	log.Printf("Guest join room_id: %s\n", room_id)

	// 2. Handle if room_id is at this instance
	// Find owner conn
	// Enter guest-owner reader
	// - Read from guest and write to owner
	// - Write to guest will be triggered by owner
	if ConnectionCache.checkRoom(room_id) {
		// 2.1 Save guest conn for owner
		err = ConnectionCache.addTarget(room_id, conn)
		if err != nil {
			return err
		}

		// 2.2 Remove guest from cache when connection closes
		handleClose := conn.CloseHandler()
		conn.SetCloseHandler(func(code int, text string) error {
			ConnectionCache.removeTarget(room_id)
			return handleClose(code, text)
		})

		// 2.3 Enter guest-owner reader
		Proxy_target_owner(room_id)

		return nil
	}

	// 3. Handle if room_id is at a different instance

	// 3.1 Find target instance IP
	owner_ip, fatal, err := Fetch_ip_from_room_id(room_id)
	if err != nil {
		if fatal == false {
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, "invalid room_id, shutdown client-guest connection."))
			return nil
		}
		return err
	}
	log.Printf("Owner instance ip: %s\n", owner_ip)

	// 3.2 Establish proxy of guest-aisle and feed AISLE_KEY and room_id
	aisle_url := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", owner_ip, port), Path: "/"}
	aisle_header := http.Header{
		"Sec-Websocket-Protocol": {"aisle"},
	}
	aisle_conn, resp, err := websocket.DefaultDialer.Dial(aisle_url.String(), aisle_header)
	if err != nil {
		err = errors.New(fmt.Sprintf("Aisle handshake failed with status %d and error: %s", resp.StatusCode, err.Error()))
		log.Println(err.Error())
		return err
	}
	if err = Supply(aisle_conn, "Authorization", os.Getenv("AISLE_KEY")); err != nil {
		return err
	}
	if err = Supply(aisle_conn, "Room-Id", room_id); err != nil {
		return err
	}

	// 3.3 Closes guest-aisle when connection client-guest closes
	guestDefaultClose := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		aisle_conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "client-guest is down, cascading to guest-aisle."))
		return guestDefaultClose(code, text)
	})

	// 3.4 Closes guest-client when connection aisle-guest closes
	aisleDefaultClose := aisle_conn.CloseHandler()
	aisle_conn.SetCloseHandler(func(code int, text string) error {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "aisle-guest is down, cascading to guest-client."))
		return aisleDefaultClose(code, text)
	})

	// 3.5 Enter guest-asile reader
	// Go routine
	// - Read from guest and write to aisle
	// - Read from aisle and write to guest
	// TODO: study if there's a better way
	go proxy_guest_aisle(conn, aisle_conn)
	proxy_aisle_guest(aisle_conn, conn)

	return nil
}

func proxy_guest_aisle(guest *websocket.Conn, aisle *websocket.Conn) error {
	for {
		messageType, body, err := guest.ReadMessage()
		if err != nil {
			err = errors.New("Proxy guest to aisle error while reading guest: " + err.Error())
			log.Println(err.Error())
			return err
		}

		err = aisle.WriteMessage(messageType, body)
		if err != nil {
			err = errors.New("Proxy guest to aisle error while writing aisle: " + err.Error())
			log.Println(err.Error())
			return err
		}
	}
}

func proxy_aisle_guest(aisle *websocket.Conn, guest *websocket.Conn) error {
	for {
		messageType, body, err := aisle.ReadMessage()
		if err != nil {
			err = errors.New("Proxy aisle to guest error while reading aisle: " + err.Error())
			log.Println(err.Error())
			return err
		}

		err = guest.WriteMessage(messageType, body)
		if err != nil {
			err = errors.New("Proxy aisle to guest error while writing guest: " + err.Error())
			log.Println(err.Error())
			return err
		}
	}
}
