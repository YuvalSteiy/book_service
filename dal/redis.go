package dal

import (
	"github.com/go-redis/redis"
)

const ADDR = "localhost:6379"

type Redis struct {
	Client *redis.Client
}

func (r *Redis) InitClient() {
	client := redis.NewClient(&redis.Options{
		Addr: ADDR,
	})
	r.Client = client
	return
}
func (r *Redis) AddActivity(username string, req string) {
	if !FieldExist(username){
		return
	}
	checkExist := r.Client.HExists(username, "path1").Val()
	if !checkExist { //no record exists for user
		r.CreateNewClientActivity(username, req)
	} else { //record already exists for client
		r.UpdateClientActivity(username, req)
	}

}
func (r *Redis) CreateNewClientActivity(username string, req string) {
	m := make(map[string]interface{})
	m["path1"] = req
	m["path2"] = ""
	m["path3"] = ""
	r.Client.HMSet(username, m)
}
func (r *Redis) UpdateClientActivity(username string, req string) {
	currRecord := r.Client.HMGet(username, "path1", "path2", "path3").Val()
	m := make(map[string]interface{})
	m["path1"] = req
	m["path2"] = currRecord[0]
	m["path3"] = currRecord[1]
	r.Client.HMSet(username, m)
}

func (r *Redis) GetUserActivity(username string) ([]interface{}) {
	if !FieldExist(username){
		return nil
	}
	if !r.Client.HExists(username,"path1").Val(){
		return nil
	}
	userData := r.Client.HMGet(username, "path1", "path2", "path3").Val()
	return userData
}
