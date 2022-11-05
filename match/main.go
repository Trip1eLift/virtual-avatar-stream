package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

func handleOwner(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "owner" {
		return nil
	}
	// TODO: check database to find unique room_id
	// TODO: register room_id with IP to the database
	message := map[string]interface{}{
		"key":   "room_id",
		"value": 1,
	}
	if payload, err := json.Marshal(message); err != nil {
		return err
	} else {
		conn.WriteMessage(websocket.TextMessage, payload)

		// save owner conn for (aisle/guest)
		// Enter owner reader
		// Dynamically find (aisle/guest) conn
		// - Read from owner and write to (aisle/guest)
		// - Write to owner will be triggered by (aisle/guest)
	}
	return nil
}

func handleGuest(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "guest" {
		return nil
	}
	room_id, _ := strconv.Atoi(request.Header["Sec-Websocket-Protocol"][1])
	log.Printf("Guest target room_id: %d\n", room_id)

	// TODO: locate instance IP from database
	// TODO: check if room_id is not in self
	owner_ip := os.Getenv("SIBILING_IP")
	aisle_url := url.URL{Scheme: "ws", Host: owner_ip, Path: "/"}

	aisle_conn, resp, err := websocket.DefaultDialer.Dial(aisle_url.String(), nil)
	// TODO: Add protocols of {"aisle", string(room_id), os.Getenv("AISLE_KEY")}

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

func handleAisle(conn *websocket.Conn, request *http.Request) error {
	if request.Header["Sec-Websocket-Protocol"][0] != "aisle" {
		return nil
	}

	// Find owner conn
	// Enter aisle-owner reader
	// - Read from aisle and write to owner
	// - Write to aisle will be triggered by owner

	return nil
}

func wsEndpoint(write http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client Successfully Connected...")

	if err := handleOwner(ws, request); err != nil {
		log.Printf("Owner error: %s\n", err.Error())
	}
	if err := handleGuest(ws, request); err != nil {
		log.Printf("Guest error: %s\n", err.Error())
	}
	reader(ws)
}

func wsHeaderCheck(request *http.Request) bool {
	if request.Header["Connection"] == nil || request.Header["Connection"][0] != "Upgrade" {
		log.Printf("Connection is not [Upgrade] but %s.\n", request.Header["Connection"])
		return false
	}
	if request.Header["Upgrade"] == nil || request.Header["Upgrade"][0] != "websocket" {
		log.Printf("Upgrade is not [websocket] but %s.\n", request.Header["Upgrade"])
		return false
	}
	if request.Header["Sec-Websocket-Key"] == nil {
		log.Println("Sec-Websocket-Key is missing.")
		return false
	}

	// Expect protocol
	// Owner len: 1 -> {"owner",}
	// Guest len: 2 -> {"guest", "<room_id>",}
	// Aisle len: 3 -> {"guest", "<room_id>", "<aisle_key>",}
	protocol := request.Header["Sec-Websocket-Protocol"]
	if protocol == nil || len(protocol) < 1 {
		log.Println("Sec-Websocket-Protocol is missing.")
		return false
	} else if protocol[0] != "owner" && protocol[0] != "guest" && protocol[0] != "aisle" {
		log.Println("Unrecognized protocol.")
		return false
	} else if protocol[0] == "owner" {
		if len(protocol) != 1 {
			log.Println("Owner protocol length error.")
			return false
		}
		if request.Header["Origin"] == nil || request.Header["Origin"][0] != os.Getenv("ORIGIN") {
			log.Printf("Owner origin is not [%s] but %s.\n", os.Getenv("ORIGIN"), request.Header["Origin"])
			return false
		}
	} else if protocol[0] == "guest" {
		if len(protocol) != 2 {
			log.Println("Guest protocol length error.")
			return false
		}
		if request.Header["Origin"] == nil || request.Header["Origin"][0] != os.Getenv("ORIGIN") {
			log.Printf("Guest origin is not [%s] but %s.\n", os.Getenv("ORIGIN"), request.Header["Origin"])
			return false
		}
		if _, err := strconv.Atoi(protocol[1]); err != nil {
			log.Printf("Guest room_id error: %s\n", err.Error())
			return false
		}
	} else if protocol[0] == "asile" {
		if len(protocol) != 3 {
			log.Println("Aisle protocol length error.")
			return false
		}
		if _, err := strconv.Atoi(protocol[1]); err != nil {
			log.Printf("Aisle room_id error: %s\n", err.Error())
			return false
		}
		if protocol[2] != os.Getenv("AISLE_KEY") {
			log.Println("Aisle key error.")
			return false
		}
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
	log.Printf("Listening on %s:%s\n", ip, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
