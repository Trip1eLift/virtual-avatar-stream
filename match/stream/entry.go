package stream

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

const ip = "0.0.0.0"
const port = "5000"

// TODO: make global functions class like
// TODO: make room_id expired after 24 hrs

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     HeaderVerify,
	Subprotocols:    []string{"owner", "guest", "aisle"},
}

func wsEndpoint(write http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Println("Client Successfully Connected...")

	if err := HandleOwner(ws, request); err != nil {
		log.Printf("Owner error")
	}
	if err := HandleGuest(ws, request, port); err != nil {
		log.Printf("Guest error")
	}
	if err := HandleAisle(ws, request); err != nil {
		log.Printf("Aisle error")
	}

	if err := ws.Close(); err != nil {
		log.Printf("Client close error: %s\n", err.Error())
	}
}

func Entry() {
	SetIp(os.Getenv("SELF_IP"))

	http.HandleFunc("/", wsEndpoint)

	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Host Address:", request.Host)
		log.Println("Remote Address:", request.RemoteAddr)
		nextRoom, err := Fetch_unique_room_id()
		if err != nil {
			panic(errors.New("Health check failed with error: " + err.Error()))
		}
		message := fmt.Sprintf("Healthy. %s\n", nextRoom)
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/database", func(write http.ResponseWriter, _ *http.Request) {
		if reply, err := Health_database(); err != nil {
			fmt.Fprintf(write, err.Error())
		} else {
			fmt.Fprintf(write, reply)
		}
	})

	log.Printf("Listening on %s:%s\n", ip, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
