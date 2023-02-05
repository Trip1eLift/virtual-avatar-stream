package stream

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// TODO: Implement logger info, debug, error so it doesn't spam message on the server

const BROADCAST_IP = "0.0.0.0"

var BROADCAST_PORT = os.Getenv("PORT")

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

	if _, err := IP.getIp(); err != nil {
		log.Printf("Owner error")
		ws.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "Owner error: "+err.Error()))
	}

	if err := HandleOwner(ws, request); err != nil {
		log.Printf("Owner error")
		ws.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "Owner error: "+err.Error()))
	}
	if err := HandleGuest(ws, request); err != nil {
		log.Printf("Guest error")
		ws.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "Guest error: "+err.Error()))
	}
	if err := HandleAisle(ws, request); err != nil {
		log.Printf("Aisle error")
		ws.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "Aisle error: "+err.Error()))
	}

	if err := ws.Close(); err != nil {
		log.Printf("Client close error: %s\n", err.Error())
	}
}

func Start() {
	if err := DBW.initializeRetry(); err != nil {
		// Two backends might try to populate table at the same time, one of them will cause an error
		log.Println("Unable to initialize database.")
	}

	if os.Getenv("SELF_IP") != "" {
		// SELF_IP should only be set in local
		IP.setIp(os.Getenv("SELF_IP"))
	}

	http.HandleFunc("/", wsEndpoint)

	http.HandleFunc("/health", func(write http.ResponseWriter, _ *http.Request) {
		self_ip, _ := IP.getIp()
		message := fmt.Sprintf("Healthy: private self IP: %s.\n", self_ip)
		//log.Print(message) // This endpoint is called too often
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/health-proxy", func(write http.ResponseWriter, _ *http.Request) {
		self_ip, _ := IP.getIp()

		if ENV == "cheap" {
			target_ip, _ := DBW.fetch_an_non_self_ip(self_ip)
			message := fmt.Sprintf("Proxy Healthy: self IP: %s target proxy IP: %s.", self_ip, target_ip)
			log.Print(message)
			fmt.Fprintf(write, message)
			return
		}

		target_ip, err := DBW.fetch_an_non_self_ip(self_ip)
		if err != nil {
			fmt.Fprintf(write, fmt.Sprintf("Database fetch error\n%s\n", err.Error()))
			return
		}

		res, err := HTTPGet(fmt.Sprintf("http://%s/health", target_ip))
		if err != nil {
			fmt.Fprintf(write, fmt.Sprintf("Proxy healthcheck error\n%s\n", err.Error()))
			return
		}

		message := fmt.Sprintf("Proxy Healthy: self IP: %s target proxy IP: %s.\nProxy health message:\n\t%s\n", self_ip, target_ip, res)
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	// ALB blocks this endpoint from public with 418
	http.HandleFunc("/health-internal", func(write http.ResponseWriter, request *http.Request) {
		IP.setIp(request.Host)
		message := fmt.Sprintf("internal health check.")
		//log.Print(message) // This endpoint is called too often
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/health-database", func(write http.ResponseWriter, _ *http.Request) {
		if reply, err := DBW.health_database(); err != nil {
			fmt.Fprintf(write, err.Error())
		} else {
			fmt.Fprintf(write, reply)
		}
	})

	log.Printf("Listening on %s:%s\n", BROADCAST_IP, BROADCAST_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", BROADCAST_IP, BROADCAST_PORT), nil))
}
