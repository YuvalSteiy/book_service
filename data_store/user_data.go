package data_store

import "sync"

var userDater UserDater
var userDaterOnce sync.Once

type UserDater interface {
	AddActivity(username string, req string)
	GetUserActivity(username string) []interface{}
}

func NewUserDater() UserDater {
	userDaterOnce.Do(func() {
		userDater = CreateRedis()
	})

	return userDater
}
