package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func HeaderVerify(request *http.Request) bool {
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
