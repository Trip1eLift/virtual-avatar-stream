package stream

import (
	"errors"
	"log"
	"sync"
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
