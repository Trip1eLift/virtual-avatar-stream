package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleAisle(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "aisle" {
		return nil
	}

	// Find owner conn
	// Enter aisle-owner reader
	// - Read from aisle and write to owner
	// - Write to aisle will be triggered by owner

	return nil
}
