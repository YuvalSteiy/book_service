package main

import (
	"github.com/gin-gonic/gin"
	"test/service"
)
func main() {
	router := gin.Default()
	service.ConfigRoutes(router)
	router.Run(":3005")
}
