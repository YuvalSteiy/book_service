package main

import (
	"github.com/YuvalSteiy/book_service/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	service.ConfigRoutes(router)
	router.Run(":3005")
}
