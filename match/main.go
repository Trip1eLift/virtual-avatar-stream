package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

const ip = "0.0.0.0"
const port = "5000"

var DB_USER = os.Getenv("DB_USER")
var DB_HOST = os.Getenv("DB_HOST")
var DB_PASS = os.Getenv("DB_PASS")
var DB_PORT = os.Getenv("DB_PORT")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     wsHeaderCheck,
	Subprotocols:    []string{"owner", "guest", "aisle"},
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

func setOwnerRoomId(conn *websocket.Conn, request *http.Request) error {
	log.Println("HERE", request.Header["Sec-Websocket-Protocol"])
	if request.Header["Sec-Websocket-Protocol"][0] == "owner" {
		// TODO: check database to find unique room_id
		message := map[string]interface{}{
			"key":   "room_id",
			"value": 1,
		}
		if payload, err := json.Marshal(message); err != nil {
			return err
		} else {
			conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
	return nil
}

func wsEndpoint(write http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client Successfully Connected...")

	setOwnerRoomId(ws, request)
	reader(ws)
}

func wsHeaderCheck(request *http.Request) bool {
	if request.Header["Connection"] == nil || request.Header["Connection"][0] != "Upgrade" {
		log.Println(fmt.Sprintf("Connection is not [Upgrade] but %s.", request.Header["Connection"]))
		return false
	}
	if request.Header["Upgrade"] == nil || request.Header["Upgrade"][0] != "websocket" {
		log.Println(fmt.Sprintf("Upgrade is not [websocket] but %s.", request.Header["Upgrade"]))
		return false
	}
	if request.Header["Sec-Websocket-Key"] == nil {
		log.Println("Sec-Websocket-Key is missing.")
		return false
	}
	if request.Header["Sec-Websocket-Protocol"] == nil || len(request.Header["Sec-Websocket-Protocol"]) <= 0 {
		log.Println("Sec-Websocket-Protocol is missing.")
		return false
	} else if (request.Header["Sec-Websocket-Protocol"][0] == "owner" ||
		request.Header["Sec-Websocket-Protocol"][0] == "guest") &&
		(request.Header["Origin"] == nil ||
			request.Header["Origin"][0] != os.Getenv("ORIGIN")) {
		log.Println(fmt.Sprintf("Origin is not [%s] but %s.", os.Getenv("ORIGIN"), request.Header["Origin"]))
		return false
	} else if request.Header["Sec-Websocket-Protocol"][0] == "aisle" &&
		(len(request.Header["Sec-Websocket-Protocol"]) != 2 ||
			request.Header["Sec-Websocket-Protocol"][1] != os.Getenv("AISLE_KEY")) {
		log.Println(fmt.Sprintf("AISLE protocol failure."))
		return false
	}
	return true
}

func main() {
	http.HandleFunc("/", wsEndpoint)
	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Host Address:", request.Host)
		log.Println("Remote Address:", request.RemoteAddr)
		log.Println("Healthy.")
		fmt.Fprintf(write, "Healthy.\n")
	})
	log.Println(fmt.Sprintf("Listening on %s:%s", ip, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
