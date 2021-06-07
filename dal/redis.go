package dal

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Redis struct{
	Client *redis.Client
}

func (r *Redis)InitClient(){
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r.Client = client
	return
}
func (r *Redis)CacheAddActivity(username string, req string)error{
	checkExist, err := r.Client.HExists(username,"path1").Result()
	if err != nil{
		return err
	}
	if !checkExist {//no record exists for user
		r.CacheCreateNewClientActivity(username,req)
	}else { //record already exists for client
		r.CacheUpdateClientActivity(username,req)
	}
	return nil
}
func (r *Redis)CacheCreateNewClientActivity(username string, req string){
	r.Client.HMSet(username,map[string]interface{}{
		"path1": req,
		"path2": "",
		"path3": "",
	})
}
func (r *Redis)CacheUpdateClientActivity(username string, req string){
	currRecord,err := r.Client.HMGet(username,"path1","path2","path3").Result()
	fmt.Println(currRecord)
	if err !=nil{
		panic(err)
	}
	r.Client.HMSet(username,map[string]interface{}{
		"path1": req,
		"path2": currRecord[0],
		"path3": currRecord[1],
	})
}

func (r *Redis) CacheGetUserActivity(username string) ([]interface{},error){
	userData,err := r.Client.HMGet(username, "path1","path2","path3").Result()
	if err != nil {
		return nil, err
	}
	return userData,nil
}
