package main

import (
	"github.com/gin-gonic/gin"
	"test/bl"
)

func main() {
	//dal.ReddisAddRecord("yuval", "this/is/a/path")
	router := gin.Default()
	configRoutes(router)
	router.Run(":3005")

}
func configRoutes(router *gin.Engine){
	router.GET("/book/:id", bl.GetBookByID)
	router.PUT("/book", bl.PutBook)
	router.POST("/book/:id", bl.UpdateBook )
	router.DELETE("/book/:id", bl.DeleteBook )
	router.GET("/search" , bl.SearchBook )
	router.GET("/store", bl.GetStoreData)
	router.GET("/activity", bl.GetUserData)
}