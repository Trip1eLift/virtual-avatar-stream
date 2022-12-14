package stream

import (
	"log"
	"net/http"
	"os"
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

	if request.Header["Sec-Websocket-Protocol"] == nil || len(request.Header["Sec-Websocket-Protocol"]) != 1 {
		log.Println("Sec-Websocket-Protocol is missing.")
		return false
	}
	protocol := request.Header["Sec-Websocket-Protocol"][0]
	if protocol != "owner" && protocol != "guest" && protocol != "aisle" {
		log.Println("Unrecognized protocol.")
		return false
	}
	if (protocol == "owner" || protocol == "guest") &&
		(request.Header["Origin"] == nil || request.Header["Origin"][0] != os.Getenv("ORIGIN")) {
		log.Printf("Guest/Owner origin is not [%s] but %s.\n", os.Getenv("ORIGIN"), request.Header["Origin"])
		return false
	}

	return true
}
