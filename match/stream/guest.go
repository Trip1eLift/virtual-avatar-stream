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
	// 1. Establish client-guest connection and retrieve room_id as int
	if request.Header["Sec-Websocket-Protocol"][0] != "guest" {
		return nil
	}
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
		// 2.1 Enter guest-owner reader
		err = ConnectionCache.addTarget(room_id_str, conn)
		if err != nil {
			return err
		}
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
	if err = Supply(conn, "Authorization", os.Getenv("AISLE_KEY")); err != nil {
		return err
	}
	if err = Supply(conn, "room_id", room_id_str); err != nil {
		return err
	}

	// 3.3 Enter guest-asile reader
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
			err = errors.New("Proxy guest to aisle error, reading guest: " + err.Error())
			log.Println(err.Error())
			return err
		}

		err = aisle.WriteMessage(messageType, body)
		if err != nil {
			err = errors.New("Proxy guest to aisle error, writing aisle: " + err.Error())
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
			err = errors.New("Proxy aisle to guest error, reading aisle: " + err.Error())
			log.Println(err.Error())
			return err
		}

		err = guest.WriteMessage(messageType, body)
		if err != nil {
			err = errors.New("Proxy aisle to guest error, writing guest: " + err.Error())
			log.Println(err.Error())
			return err
		}
	}
}
