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
	CheckOrigin:     wsHeaderCheck,
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
	response_header := http.Header{
		"Status": {"Connection-Established"},
		//"Room":   {"00001"},
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
	if request.Header["Origin"] == nil || request.Header["Origin"][0] != os.Getenv("ORIGIN") {
		log.Println(fmt.Sprintf("Origin is not [%s] but %s.", os.Getenv("ORIGIN"), request.Header["Origin"]))
		return false
	}
	if request.Header["Sec-Websocket-Protocol"] == nil || len(request.Header["Sec-Websocket-Protocol"]) <= 0 {
		log.Println("Sec-Websocket-Protocol is missing.")
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
		log.Println("Host Address:", request.Host)
		log.Println("Remote Address:", request.RemoteAddr)
		log.Println("Healthy.")
		fmt.Fprintf(write, "Healthy.\n")
	})
	log.Println(fmt.Sprintf("Listening on %s:%s", ip, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
