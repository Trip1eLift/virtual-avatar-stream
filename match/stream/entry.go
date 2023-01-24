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

	// TODO: close connection if self IP was not set; Self ip is always set, skip this for now

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
	if err := DB.initializeRetry(); err != nil {
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
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/health-proxy", func(write http.ResponseWriter, _ *http.Request) {
		self_ip, _ := IP.getIp()
		target_ip, err := DB.fetch_an_non_self_ip(self_ip)
		if err != nil {
			fmt.Fprintf(write, fmt.Sprintf("Database fetch error\n%s\n", err.Error()))
			return
		}

		res, err := HTTPGet(fmt.Sprintf("http://%s:%s/health", target_ip, port))
		if err != nil {
			fmt.Fprintf(write, fmt.Sprintf("Proxy healthcheck error\n%s\n", err.Error()))
			return
		}

		// TODO PRIO: health-proxy not working on aws:
		//						Proxy healthcheck error
		// 						Execute get reqeust error: Get "http://10.0.8.198:5000:5000/health": dial tcp: lookup 10.0.8.198:5000: no such host

		message := fmt.Sprintf("Proxy Healthy: self IP: %s target proxy IP: %s.\nProxy health message:\n\t%s\n", self_ip, target_ip, res)
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	// This endpoint should not be hit by public
	// TODO: check if ALB can block this
	http.HandleFunc("/health-internal", func(write http.ResponseWriter, request *http.Request) {
		IP.setIp(request.Host)
		message := fmt.Sprintf("internal health check.")
		log.Print(message)
		fmt.Fprintf(write, message)
	})

	http.HandleFunc("/health-database", func(write http.ResponseWriter, _ *http.Request) {
		if reply, err := DB.health_database(); err != nil {
			fmt.Fprintf(write, err.Error())
		} else {
			fmt.Fprintf(write, reply)
		}
	})

	log.Printf("Listening on %s:%s\n", ip, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
