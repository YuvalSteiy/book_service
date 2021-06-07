package service

import (
	"github.com/gin-gonic/gin"
	"test/bl"
)

func ConfigRoutes(router *gin.Engine){
	router.GET("/book/:id", bl.GetBookByID)
	router.PUT("/book", bl.PutBook)
	router.POST("/book/:id", bl.UpdateBook )
	router.DELETE("/book/:id", bl.DeleteBook )
	router.GET("/search" , bl.SearchBook )
	router.GET("/store", bl.GetStoreData)
	router.GET("/activity", bl.GetUserData)
}