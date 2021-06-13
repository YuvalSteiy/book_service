package main

import (
	"github.com/YuvalSteiy/book_service/service"
	"github.com/gin-gonic/gin"
)

const SERVER_PORT = ":3005"

func main() {
	server := gin.Default()
	service.ConfigRoutes(server)
	server.Run(SERVER_PORT)
}
