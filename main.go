package main

import (
	"fmt"
	"github.com/YuvalSteiy/book_service/service"
	"github.com/gin-gonic/gin"
)

const SERVER_PORT = "3005"

func main() {
	server := gin.Default()
	service.ConfigRoutes(server)
	err := service.InitDataStore()
	if err != nil {
		panic(err)
	}
	server.Run(fmt.Sprintf(":%s", SERVER_PORT))
}
