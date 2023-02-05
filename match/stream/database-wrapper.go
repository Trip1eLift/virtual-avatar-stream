// To save cost in cheap deploy mode, it doesn't create a database at all.
// It uses memory and generated room_id with target private ip.

// Replace function list:
// 1. save_room_id_with_ip(room_id string, ip string) error
//    room_id: 10.0.2.51:5000:3, ip: ""
//    return nil
//
// 2. remove_room_id(room_id string) error
//    return nil
//
// 3. fetch_ip_from_room_id(room_id string) (string, bool, error)
//    room_id: 10.0.2.51:5000:3
//    detach 10.0.2.51:5000 from room_id
//    return 10.0.2.51:5000, false, nil
//
// 4. fetch_unique_room_id() (string, error)
//    keep a local counter with mutex lock
//    return 3, nil
//
// 5. health_database() (string, error)
//    return "cheap mode no database", nil
//
// 6. fetch_an_non_self_ip(self_ip string) (string, error)
//    return "cheap mode no target ips", nil
//    need to make code changes after DB.fetch_an_non_self_ip
//
// 7. initializeRetry() error
//    return nil
//
// 8. initialize() error
//    return nil
//
// 9. cleanup(self_ip string) error
//    return nil

package stream

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
)

var ENV = os.Getenv("environment")

type DatabaseWrapper struct {
	room_number uint32
	lock_number sync.RWMutex
}

func (d *DatabaseWrapper) save_room_id_with_ip(room_id string, ip string) error {
	if ENV != "cheap" {
		return DB.save_room_id_with_ip(room_id, ip)
	}
	return nil
}

func (d *DatabaseWrapper) remove_room_id(room_id string) error {
	if ENV != "cheap" {
		return DB.remove_room_id(room_id)
	}
	return nil
}

func (d *DatabaseWrapper) fetch_ip_from_room_id(room_id string) (string, bool, error) {
	if ENV != "cheap" {
		return DB.fetch_ip_from_room_id(room_id)
	}

	if count := strings.Count(room_id, ":"); count != 2 {
		err := errors.New("[cheap mode] room_id does not contain 2 ':'.")
		return "", false, err
	}

	i1 := strings.Index(room_id, ":")
	room_id_sub := room_id[i1+1:]
	i2 := strings.Index(room_id_sub, ":")

	target_ip := room_id[:i1+1+i2]
	return target_ip, false, nil
}

func (d *DatabaseWrapper) fetch_unique_room_id() (string, error) {
	if ENV != "cheap" {
		return DB.fetch_unique_room_id()
	}

	d.lock_number.Lock()
	room_number := d.room_number
	d.room_number++
	d.lock_number.Unlock()

	ip, err := IP.getIp()
	if err != nil {
		return "", err
	}
	room_id := ip + ":" + strconv.FormatInt(int64(room_number), 10)

	return room_id, nil
}

func (d *DatabaseWrapper) health_database() (string, error) {
	if ENV != "cheap" {
		return DB.health_database()
	}
	return "cheap mode no database", nil
}

func (d *DatabaseWrapper) fetch_an_non_self_ip(self_ip string) (string, error) {
	if ENV != "cheap" {
		return DB.fetch_an_non_self_ip(self_ip)
	}
	return "cheap mode no target ips", nil
}

func (d *DatabaseWrapper) initializeRetry() error {
	if ENV != "cheap" {
		return DB.initializeRetry()
	}
	return nil
}

func (d *DatabaseWrapper) initialize() error {
	if ENV != "cheap" {
		return DB.initialize()
	}
	return nil
}

func (d *DatabaseWrapper) cleanup(self_ip string) error {
	if ENV != "cheap" {
		return DB.cleanup(self_ip)
	}
	return nil
}

var DBW = DatabaseWrapper{
	room_number: 1,
}
