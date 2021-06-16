package data_store

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const ADDR = "localhost:6379"

type redisUserData struct {
	client *redis.Client
}

func initRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: ADDR})
}

func newRedisUserData() *redisUserData {
	client := initRedisClient()
	if client == nil {
		return nil
	}

	return &redisUserData{client: client}
}

func (r *redisUserData) AddActivity(username string, req string) error {
	checkExist := r.client.HExists(username, "path1").Val()
	var err error
	if !checkExist {
		err = r.createNewClientActivity(username, req)
	} else {
		err = r.updateClientActivity(username, req)
	}
	return err
}

func (r *redisUserData) createNewClientActivity(username string, req string) error {
	m := map[string]interface{}{
		"path1": req,
		"path2": "",
		"path3": "",
	}

	_, err := r.client.HMSet(username, m).Result()
	return err
}

func (r *redisUserData) updateClientActivity(username string, req string) error {
	currRecord, err := r.client.HMGet(username, "path1", "path2", "path3").Result()
	if err != nil {
		return err
	}

	m := map[string]interface{}{
		"path1": req,
		"path2": currRecord[0],
		"path3": currRecord[1],
	}

	_, err = r.client.HMSet(username, m).Result()
	return err
}

func (r *redisUserData) GetUserActivity(username string) ([]interface{}, error) {
	exist, err := r.client.HExists(username, "path1").Result()
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("No Activity For This User")
	}

	userData, err := r.client.HMGet(username, "path1", "path2", "path3").Result()
	if err != nil {
		return nil, err
	}

	return userData, nil
}
