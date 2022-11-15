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
	wait   *sync.WaitGroup
}

type Connections struct {
	index map[string]*ConnectionRoom
	mutex sync.RWMutex
}

func (c *Connections) addRoom(room_id string, owner *websocket.Conn) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == true {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot add room_id: %s because it is already in index.", room_id))
		log.Println(err.Error())
		return err
	}
	c.index[room_id] = &ConnectionRoom{
		owner:  owner,
		target: nil,
		wait:   &sync.WaitGroup{},
	}
	c.index[room_id].wait.Add(1)
	c.mutex.Unlock()
	return nil
}

func (c *Connections) removeRoom(room_id string) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot remove room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return err
	}
	if targetConn := c.index[room_id].target; targetConn != nil {
		// close target connection gracefully
		targetConn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "client-owner is down, cascading to owner-target."))
	}
	delete(c.index, room_id)
	c.mutex.Unlock()
	return nil
}

func (c *Connections) waitRoom(room_id string) (bool, error) {
	terminate := false
	c.mutex.RLock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.RUnlock()
		err := errors.New(fmt.Sprintf("Cannot wait on room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return terminate, err
	}
	owner := c.index[room_id].owner
	wait := c.index[room_id].wait
	c.mutex.RUnlock()

	// Catch if owner close connection before any target appears
	quit := make(chan bool)
	connected := false
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				if terminate {
					break // force loop to be alive after terminate and before quit
				}
				messageType, body, err := owner.ReadMessage()
				if err != nil {
					wait.Done()
					terminate = true
				}

				// Catch the first read after owner-target connection is established
				if connected {
					target, err := c.getTarget(room_id)
					if err == nil {
						err = target.WriteMessage(messageType, body)
						if err != nil {
							log.Printf("Catch first read error: %s\n", err.Error())
						}
					}
				}
			}
		}
	}()

	wait.Wait() // wait for room_id be signal to start
	connected = true
	quit <- true
	return terminate, nil
}

func (c *Connections) getTarget(room_id string) (*websocket.Conn, error) {
	c.mutex.RLock()
	room := c.index[room_id]
	if room == nil || room.target == nil {
		c.mutex.RUnlock()
		err := errors.New(fmt.Sprintf("Cannot get target from room_id: %s", room_id))
		log.Println(err.Error())
		return nil, err
	}
	c.mutex.RUnlock()
	return room.target, nil
}

// The return bool indicates if the error is fatal to owner
func (c *Connections) getRoom(room_id string) (*websocket.Conn, *websocket.Conn, bool, error) {
	c.mutex.RLock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.RUnlock()
		err := errors.New(fmt.Sprintf("Cannot get room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return nil, nil, true, err
	}
	owner := c.index[room_id].owner
	target := c.index[room_id].target
	c.mutex.RUnlock()
	if target == nil {
		err := errors.New(fmt.Sprintf("Cannot get room_id: %s because target was not assigned.", room_id))
		log.Println(err.Error())
		return nil, nil, false, err
	}
	return owner, target, false, nil
}

func (c *Connections) checkRoom(room_id string) bool {
	c.mutex.RLock()
	_, found := c.index[room_id]
	c.mutex.RUnlock()
	return found
}

func (c *Connections) addTarget(room_id string, target *websocket.Conn) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot add target to room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return err
	}
	c.index[room_id].target = target
	c.index[room_id].wait.Done()
	c.mutex.Unlock()
	return nil
}

func (c *Connections) removeTarget(room_id string) error {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot remove target from room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return err
	}
	c.index[room_id].target = nil
	c.index[room_id].wait.Add(1)
	c.mutex.Unlock()
	return nil
}

// Declare Global variable
var ConnectionCache = Connections{index: make(map[string]*ConnectionRoom)}
