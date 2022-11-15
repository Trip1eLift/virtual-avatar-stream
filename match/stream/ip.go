package stream

import (
	"errors"
	"log"
	"sync"
)

var IP = ""
var IPmutex sync.RWMutex

func SetIp(ip string) {
	IPmutex.RLock()
	if IP != "" {
		IPmutex.RUnlock()
		return
	}
	IPmutex.RUnlock()
	IPmutex.Lock()
	IP = ip
	log.Printf("Set self IP to be: %s", IP)
	IPmutex.Unlock()
}

func GetIp() (string, error) {
	IPmutex.RLock()
	if IP == "" {
		IPmutex.RUnlock()
		err := errors.New("Cannot retrieve self IP because it was not set.")
		log.Println(err)
		return "", err
	}
	ip := IP
	IPmutex.RUnlock()
	return ip, nil
}
