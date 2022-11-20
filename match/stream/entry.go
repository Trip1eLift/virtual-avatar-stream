package stream

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

const ip = "0.0.0.0"
const port = "5000"

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

	// TODO: close connection if self IP was not set

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

func Start() {
	if os.Getenv("SELF_IP") != "" {
		// SELF_IP should only be set in local
		IP.setIp(os.Getenv("SELF_IP"))
	}

	http.HandleFunc("/", wsEndpoint)

	// TODO: make a proxy-health endpoint to test client->fargate->fargate connection

	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		self_ip, _ := IP.getIp()
		message := fmt.Sprintf("Healthy. private self IP: %s.\n", self_ip)
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/internal-health", func(write http.ResponseWriter, request *http.Request) {
		IP.setIp(request.Host)
		message := fmt.Sprintf("internal health check.")
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/database", func(write http.ResponseWriter, _ *http.Request) {
		if reply, err := DB.health_database(); err != nil {
			fmt.Fprintf(write, err.Error())
		} else {
			fmt.Fprintf(write, reply)
		}
	})

	log.Printf("Listening on %s:%s\n", ip, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
