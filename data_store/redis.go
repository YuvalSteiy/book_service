package data_store

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const ADDR = "localhost:6379"

type redisClient struct {
	client *redis.Client
}

func initRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: ADDR})
}

func NewRedisClient() *redisClient {
	client := initRedisClient()
	if client == nil {
		return nil
	}

	return &redisClient{client: client}
}

func (r *redisClient) AddActivity(username string, req string) error {
	checkExist := r.client.HExists(username, "path1").Val()
	if !checkExist {
		err := r.createNewClientActivity(username, req)
		if err != nil {
			return err
		}
	} else {
		err := r.updateClientActivity(username, req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *redisClient) createNewClientActivity(username string, req string) error {
	m := map[string]interface{}{
		"path1": req,
		"path2": "",
		"path3": "",
	}

	_, err := r.client.HMSet(username, m).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClient) updateClientActivity(username string, req string) error {
	currRecord, err := r.client.HMGet(username, "path1", "path2", "path3").Result()
	if err != nil {
		return err
	}

	m := map[string]interface{}{
		"path1": req,
		"path2": currRecord[0],
		"path3": currRecord[1],
	}

	r.client.HMSet(username, m)
	return nil
}

func (r *redisClient) GetUserActivity(username string) ([]interface{}, error) {
	exist, err := r.client.HExists(username, "path1").Result()
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("No Activity For This User")
	}

	userData, err := r.client.HMGet(username, "path1", "path2", "path3").Result()
	if err != nil{
		return nil, err
	}

	return userData, nil
}
