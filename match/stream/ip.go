package stream

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Ip struct {
	ip string
	mu sync.RWMutex
}

func (i *Ip) setIp(ip string) {
	i.mu.RLock()
	if i.ip != "" {
		i.mu.RUnlock()
		return
	}
	i.mu.RUnlock()
	i.mu.Lock()
	i.ip = ip
	log.Printf("Set self IP to be: %s", i.ip)
	i.mu.Unlock()

	go reserveRoom(ip)
}

func reserveRoom(ip string) {
	// reserve temporary room (24hr) for task itself using negative room_id, this is for health-proxy
	for i := 0; i < 10; i++ {
		room_id, err := DB.fetch_unique_room_id()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		err = DB.save_room_id_with_ip("-"+room_id, ip)
		if err == nil {
			return
		}
	}
}

func (i *Ip) getIp() (string, error) {
	i.mu.RLock()
	if i.ip == "" {
		i.mu.RUnlock()
		err := errors.New("Cannot retrieve self IP because it was not set.")
		log.Println(err)
		return "", err
	}
	ip := i.ip
	i.mu.RUnlock()
	return ip, nil
}

var IP = Ip{ip: ""}
