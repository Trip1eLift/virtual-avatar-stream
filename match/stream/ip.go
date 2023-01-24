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
		// Do nothing if self_ip is set. Neither set self_ip or populate database schema
		return
	}
	i.mu.RUnlock()
	i.mu.Lock()
	i.ip = ip
	log.Printf("Set self IP to be: %s", i.ip)
	i.mu.Unlock()

	// TODO: this can be turn off after testing
	// reserve temporary room (24hr) for task itself using negative room_id, this is for health-proxy
	go func(ip string) {
		// Retry 6 times - local: 6 sec - aws: 3 min
		for i := 0; i < 6; i++ {
			room_id, err := DB.fetch_unique_room_id()
			if err != nil {
				time.Sleep(time.Duration(DB.backoff) * time.Second)
				continue
			}
			err = DB.save_room_id_with_ip("-"+room_id, ip)
			if err == nil {
				return
			}
		}
	}(ip)
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
