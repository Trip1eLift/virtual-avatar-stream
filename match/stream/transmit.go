package stream

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

type Payload struct {
	Bus string
}

type Transmit struct {
}

// Write then wait on read
// Make sure ask follows the http header format
func (t *Transmit) demand(conn *websocket.Conn, ask string) (string, error) {
	pack := Payload{Bus: ask}
	body, err := json.Marshal(pack)
	if err != nil {
		err = errors.New("Demand marshal error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}

	err = conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		err = errors.New("Demand write error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}

	_, body, err = conn.ReadMessage()
	if err != nil {
		err = errors.New("Demand read error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}

	var ans Payload
	err = json.Unmarshal(body, &ans)
	if err != nil {
		err = errors.New("Demand unmarshal error: " + err.Error())
		log.Println(err.Error())
		return "", err
	}
	return ans.Bus, nil
}

// Wait on read then write
// Make sure ask follows the http header format
func (t *Transmit) supply(conn *websocket.Conn, ask string, ans string) error {
	_, body, err := conn.ReadMessage()
	if err != nil {
		err = errors.New("Supply read error: " + err.Error())
		log.Println(err.Error())
		return err
	}

	var pack Payload
	err = json.Unmarshal(body, &pack)
	if err != nil {
		log.Println("DEBUG", ask, ans, string(body))
		err = errors.New("Supply unmarshal error: " + err.Error())
		log.Println(err.Error())
		return err
	}

	if pack.Bus != ask {
		log.Println("DEBUG", ask, ans, pack)
		err = errors.New("Supply verifying error: expected " + ask + " but recieved " + pack.Bus)
		log.Println(err.Error())
		return err
	}

	pack = Payload{Bus: ans}
	body, err = json.Marshal(pack)
	if err != nil {
		err = errors.New("Supply marshal error: " + err.Error())
		log.Println(err.Error())
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		err = errors.New("Supply write error: " + err.Error())
		log.Println(err.Error())
		return err
	}
	return nil
}

var TM = Transmit{}
