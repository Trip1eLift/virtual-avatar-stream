package stream

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func HeaderVerify(request *http.Request) bool {
	if request.Header["Connection"] == nil || strings.ToLower(request.Header["Connection"][0]) != "upgrade" {
		log.Printf("Connection is not [upgrade] but %s.\n", request.Header["Connection"])
		return false
	}
	if request.Header["Upgrade"] == nil || strings.ToLower(request.Header["Upgrade"][0]) != "websocket" {
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

func HTTPGet(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.New("New get reqeust error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New("Execute get reqeust error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read get reqeust body error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}

	return string(body), nil
}
