package service

import (
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(router *gin.Engine){
	router.GET("/book/:id", GetBookByID)
	router.PUT("/book", InsertBook)
	router.POST("/book/:id", UpdateBook)
	router.DELETE("/book/:id", DeleteBook)
	router.GET("/search" , SearchBook)
	router.GET("/store", GetStoreInfo)
	router.GET("/activity", GetUserData)
}