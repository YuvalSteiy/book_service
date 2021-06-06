package dal

import (
	"fmt"
	"github.com/go-redis/redis"
)

func rClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}
func ReddisAddRecord(username string, req string){
	client := rClient()
	checkExist, err := client.HExists(username,"path1").Result()
	if err != nil{
		panic(err)
	}
	if !checkExist {//no record exists for user
		reddisCreateNewCLientRecord(username,req,client)
	}else { //record already exists for client
		reddisUpdateClientRecord(username,req,client)
	}
	client.Close()
}
func reddisCreateNewCLientRecord(username string, req string,client *redis.Client ){
	client.HMSet(username,map[string]interface{}{
		"path1": req,
		"path2": "---",
		"path3": "---",
	})
}
func reddisUpdateClientRecord(username string, req string, client *redis.Client){
	fmt.Println("Current user record")
	currRecord,err := client.HMGet(username,"path1","path2","path3").Result()
	if err !=nil{
		panic(err)
	}
	fmt.Println(currRecord)
	client.HMSet(username,map[string]interface{}{
		"path1": req,
		"path2": currRecord[0],
		"path3": currRecord[1],
	})
}
func ReddisGetUserRecord(username string) []interface{}{
	client := rClient()
	userData,err := client.HMGet(username, "path1","path2","path3").Result()
	if err != nil{
		panic(err)
	}
	client.Close()
	return userData
}
