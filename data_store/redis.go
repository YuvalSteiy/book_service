package data_store

import (
	"github.com/go-redis/redis"
)

const ADDR = "localhost:6379"

type Redis struct {
	Client *redis.Client
}

func initRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: ADDR})
}

func CreateRedis() *Redis {
	r := Redis{Client: initRedisClient()}
	if r.Client == nil {
		return nil
	}

	return &r
}

func (r *Redis) AddActivity(username string, req string) {
	checkExist := r.Client.HExists(username, "path1").Val()
	if !checkExist {
		r.createNewClientActivity(username, req)
	} else {
		r.updateClientActivity(username, req)
	}

}

func (r *Redis) createNewClientActivity(username string, req string) {
	m := map[string]interface{}{
		"path1": req,
		"path2": "",
		"path3": "",
	}

	r.Client.HMSet(username, m)
}

func (r *Redis) updateClientActivity(username string, req string) {
	currRecord := r.Client.HMGet(username, "path1", "path2", "path3").Val()
	m := map[string]interface{}{
		"path1": req,
		"path2": currRecord[0],
		"path3": currRecord[1],
	}

	r.Client.HMSet(username, m)
}

func (r *Redis) GetUserActivity(username string) []interface{} {
	if !r.Client.HExists(username, "path1").Val() {
		return nil
	}

	userData := r.Client.HMGet(username, "path1", "path2", "path3").Val()
	return userData
}
