package stream

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func HandleGuest(conn *websocket.Conn, request *http.Request, port string) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "guest" {
		return nil
	}

	// 1. Establish client-guest connection and retrieve room_id as int
	room_id_str, err := Demand(conn, "Room-Id")
	if err != nil {
		return err
	}
	room_id, err := strconv.Atoi(room_id_str)
	if err != nil {
		err = errors.New("Guest room_id type cast error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	log.Printf("Guest join room_id: %d\n", room_id)

	// 2. Handle if room_id is at this instance
	// Find owner conn
	// Enter guest-owner reader
	// - Read from guest and write to owner
	// - Write to guest will be triggered by owner
	if ConnectionCache.checkRoom(room_id_str) {
		// 2.1 Save guest conn for owner
		err = ConnectionCache.addTarget(room_id_str, conn)
		if err != nil {
			return err
		}

		// 2.2 Remove guest from cache when connection closes
		handleClose := conn.CloseHandler()
		conn.SetCloseHandler(func(code int, text string) error {
			ConnectionCache.removeTarget(room_id_str)
			return handleClose(code, text)
		})

		// 2.3 Enter guest-owner reader
		Proxy_target_owner(room_id_str)

		return nil
	}

	// 3. Handle if room_id is at a different instance

	// 3.1 Find target instance IP
	// TODO: locate instance IP from database
	// TODO: check if room_id is not in self
	owner_ip := os.Getenv("SIBILING_IP")
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
	if err = Supply(aisle_conn, "Room-Id", room_id_str); err != nil {
		return err
	}

	// 3.3 Closes guest-aisle when connection client-guest closes
	handleClose := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("DEBUG client-guest is down, cascading to guest-aisle.")
		aisle_conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "client-guest is down, cascading to guest-aisle."))
		return handleClose(code, text)
	})

	// 3.4 Enter guest-asile reader
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
		// TODO: return nil if guest close gracefully
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
		// TODO: return nil if aisle close gracefully
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
