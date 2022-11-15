package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Trip1eLift/virtual-avatar-stream/stream"
	"github.com/gorilla/websocket"
)

const ip = "0.0.0.0"
const port = "5000"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     stream.HeaderVerify,
	Subprotocols:    []string{"owner", "guest", "aisle"},
}

func wsEndpoint(write http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Println("Client Successfully Connected...")

	if err := stream.HandleOwner(ws, request); err != nil {
		log.Printf("Owner error")
	}
	if err := stream.HandleGuest(ws, request, port); err != nil {
		log.Printf("Guest error")
	}
	if err := stream.HandleAisle(ws, request); err != nil {
		log.Printf("Aisle error")
	}

	if err := ws.Close(); err != nil {
		log.Printf("Client close error: %s\n", err.Error())
	}
}

func main() {
	stream.SetIp(os.Getenv("SELF_IP"))

	http.HandleFunc("/", wsEndpoint)
	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Host Address:", request.Host)
		log.Println("Remote Address:", request.RemoteAddr)
		//stream.Fetch_unique_room_id()
		log.Println("Healthy.")
		fmt.Fprintf(write, "Healthy.\n")
	})
	http.HandleFunc("/database", func(write http.ResponseWriter, _ *http.Request) {
		if reply, err := stream.Health_database(); err != nil {
			fmt.Fprintf(write, err.Error())
		} else {
			fmt.Fprintf(write, reply)
		}
	})
	log.Printf("Listening on %s:%s\n", ip, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
