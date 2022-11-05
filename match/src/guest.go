package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func HandleGuest(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "guest" {
		return nil
	}
	room_id, _ := strconv.Atoi(request.Header["Sec-Websocket-Protocol"][1])
	log.Printf("Guest target room_id: %d\n", room_id)

	// TODO: locate instance IP from database
	// TODO: check if room_id is not in self
	owner_ip := os.Getenv("SIBILING_IP")
	aisle_url := url.URL{Scheme: "ws", Host: owner_ip, Path: "/"}

	aisle_header := http.Header{
		"Sec-Websocket-Protocol": {"aisle", strconv.Itoa(room_id), os.Getenv("AISLE_KEY")},
	}
	aisle_conn, resp, err := websocket.DefaultDialer.Dial(aisle_url.String(), aisle_header)

	if err != nil {
		log.Printf("Aisle handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
		return err
	}
	// Enter guest-asile reader
	// Go routine
	// - Read from guest and write to aisle
	// - Read from aisle and write to guest
	aisle_conn.WriteMessage(websocket.TextMessage, []byte("Hello"))

	// TODO: handle if self is the matching instance
	// Find owner conn
	// Enter guest-owner reader
	// - Read from guest and write to owner
	// - Write to guest will be triggered by owner
	return nil
}
