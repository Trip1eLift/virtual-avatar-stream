package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const port = "5001"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func wsEndpoint(write http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = wsHeaderCheck
	response_header := http.Header{
		"Status": {"Connection-Established"},
	}
	ws, err := upgrader.Upgrade(write, request, response_header)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client Successfully Connected...")

	if close := wsHealthCheck(ws, request); close == true {
		return
	}
	reader(ws)
}

func wsHeaderCheck(request *http.Request) bool {
	if request.Header["Connection"] == nil || request.Header["Connection"][0] != "Upgrade" {
		log.Println("Header Connection is not [Upgrade] but", request.Header["Connection"])
		return false
	}
	if request.Header["Upgrade"] == nil || request.Header["Upgrade"][0] != "websocket" {
		log.Println("Header Upgrade is not [websocket] but", request.Header["Upgrade"])
		return false
	}
	if request.Header["Sec-Websocket-Key"] == nil {
		log.Println("Header Sec-Websocket-Key is missing")
		return false
	}
	return true
}

func wsHealthCheck(conn *websocket.Conn, request *http.Request) bool {
	if request.Header["Health"] == nil || request.Header["Health"][0] != "healthcheck" {
		return false
	}
	log.Println("This is a websocket healthcheck connection. Instance is healthy. Disconnecting...")
	// 8: CloseMessage, 1000: CloseNormalClosure
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Instance is healthy!")); err != nil {
		log.Println(err)
	}
	if err := conn.Close(); err != nil {
		log.Println(err)
	}
	return true
}

func main() {
	http.HandleFunc("/", wsEndpoint)
	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Healthy.")
		fmt.Fprintf(write, "Healthy")
	})
	log.Println("Listening on localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
