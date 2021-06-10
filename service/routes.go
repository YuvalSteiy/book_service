package service

import (
	"github.com/YuvalSteiy/book_service/bl"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(router *gin.Engine){
	router.GET("/book/:id", bl.GetBookByID)
	router.PUT("/book", bl.InsertBook)
	router.POST("/book/:id", bl.UpdateBook )
	router.DELETE("/book/:id", bl.DeleteBook )
	router.GET("/search" , bl.SearchBook )
	router.GET("/store", bl.GetStoreInfo)
	router.GET("/activity", bl.GetUserData)
}