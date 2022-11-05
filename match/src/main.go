package main

import (
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
	CheckOrigin:     HeaderVerify,
	Subprotocols:    []string{"owner", "guest", "aisle"},
}

func wsEndpoint(write http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client Successfully Connected...")

	if err := HandleOwner(ws, request); err != nil {
		log.Printf("Owner error: %s\n", err.Error())
	}
	if err := HandleGuest(ws, request); err != nil {
		log.Printf("Guest error: %s\n", err.Error())
	}
	if err := HandleAisle(ws, request); err != nil {
		log.Printf("Aisle error: %s\n", err.Error())
	}

	if err := ws.Close(); err != nil {
		log.Printf("Client close error: %s\n", err.Error())
	}
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
