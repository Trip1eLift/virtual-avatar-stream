package stream

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionRoom struct {
	owner  *websocket.Conn
	target *websocket.Conn // target can be either guest or aisle
}

type Connections struct {
	index map[string]*ConnectionRoom
	mutex sync.Mutex
}

func (c *Connections) addRoom(room_id string, owner *websocket.Conn) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == false {
		c.index[room_id] = &ConnectionRoom{
			owner:  owner,
			target: nil,
		}
		c.mutex.Unlock()
		return nil
	} else {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot add room_id: %s because it is already in index.", room_id))
		log.Println(err.Error())
		return err
	}
}

func (c *Connections) removeRoom(room_id string) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == true {
		delete(c.index, room_id)
		c.mutex.Unlock()
		return nil
	} else {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot remove room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return err
	}
}

func (c *Connections) addTarget(room_id string, target *websocket.Conn) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == true {
		c.index[room_id].target = target
		c.mutex.Unlock()
		return nil
	} else {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot add target to room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return err
	}
}

// Declare Global variable
var ConnectionCache Connections
