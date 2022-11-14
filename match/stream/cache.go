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
	mutex sync.Mutex
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
	delete(c.index, room_id)
	c.mutex.Unlock()
	return nil
}

func (c *Connections) waitRoom(room_id string) (bool, error) {
	terminate := false
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot wait on room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return terminate, err
	}
	owner := c.index[room_id].owner
	wait := c.index[room_id].wait
	c.mutex.Unlock()

	// Catch if owner close connection before any target appears
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				if terminate {
					break // force loop to be alive after terminate and before quit
				}
				// TODO: This will eat away one read from owner after owner-target (guest/aisle) is established
				//       Is it possible to fix it?
				if _, _, err := owner.NextReader(); err != nil {
					wait.Done()
					terminate = true
				}
			}
		}
	}()

	wait.Wait() // wait for room_id be signal to start
	quit <- true
	return terminate, nil
}

// The return bool indicates if the error is fatal to owner
func (c *Connections) getRoom(room_id string) (*websocket.Conn, *websocket.Conn, bool, error) {
	c.mutex.Lock()
	if _, found := c.index[room_id]; found == false {
		c.mutex.Unlock()
		err := errors.New(fmt.Sprintf("Cannot get room_id: %s because it does not exist in index.", room_id))
		log.Println(err.Error())
		return nil, nil, true, err
	}
	owner := c.index[room_id].owner
	target := c.index[room_id].target
	c.mutex.Unlock()
	if target == nil {
		err := errors.New(fmt.Sprintf("Cannot get room_id: %s because target was not assigned.", room_id))
		log.Println(err.Error())
		return nil, nil, false, err
	}
	return owner, target, false, nil
}

func (c *Connections) checkRoom(room_id string) bool {
	c.mutex.Lock()
	_, found := c.index[room_id]
	c.mutex.Unlock()
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
