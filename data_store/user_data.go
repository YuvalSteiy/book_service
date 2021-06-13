package data_store

<<<<<<< HEAD
import "sync"

var userDater UserDater
var userDaterOnce sync.Once
=======
var userDater UserDater
>>>>>>> 3307927fd33aca1aeb0a88015e2937965045df37

type UserDater interface {
	AddActivity(username string, req string)
	GetUserActivity(username string) []interface{}
}

func NewUserDater() UserDater {
<<<<<<< HEAD
	userDaterOnce.Do(func(){
		userDater = CreateRedis()
	})
=======
	if userDater == nil {
		userDater = CreateRedis()
	}
>>>>>>> 3307927fd33aca1aeb0a88015e2937965045df37

	return userDater
}
