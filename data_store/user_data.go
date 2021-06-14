package data_store

import "sync"

type UserDater interface {
	AddActivity(username string, req string) error
	GetUserActivity(username string) ([]interface{}, error)
}

var userDater UserDater
var userDaterOnce sync.Once

func NewUserDater() UserDater {
	userDaterOnce.Do(func() {
		userDater = newRedisUserData()
	})

	return userDater
}
