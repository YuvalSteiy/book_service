package data_store

var userDater UserDater

type UserDater interface {
	AddActivity(username string, req string)
	GetUserActivity(username string) []interface{}
}

func NewUserDater() UserDater {
	if userDater == nil {
		userDater = CreateRedis()
	}

	return userDater
}
