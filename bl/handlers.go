package bl

import (
	"encoding/json"
	"fmt"
	"github.com/YuvalSteiy/book_service/dal"
	"github.com/YuvalSteiy/book_service/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetBookByID(c *gin.Context) {
	updateUserData(c, "GET")
	db, err := dal.NewBookStorer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	id := c.Param("id")
	book, err := db.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, book)
	return
}

func InsertBook(c *gin.Context) {
	updateUserData(c, "PUT")
	db, err := dal.NewBookStorer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var book models.Book
	json.Unmarshal(jsonData, &book)
	putID, err := db.InsertBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Book inserted at ID: "+putID)
}

func UpdateBook(c *gin.Context) {
	updateUserData(c, "POST")
	db, err := dal.NewBookStorer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	bookId := c.Param("id")
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var book models.Book
	json.Unmarshal(jsonData, &book)
	err = db.UpdateBook(*book.Title, bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Successes")
}

func DeleteBook(c *gin.Context) {
	updateUserData(c, "Delete")
	db, err := dal.NewBookStorer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	id := c.Param("id")
	err = db.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Successes")
}

func SearchBook(c *gin.Context) {
	updateUserData(c, "GET")
	db, err := dal.NewBookStorer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	title := c.DefaultQuery("title", "")
	authorName := c.DefaultQuery("author_name", "")
	priceRangeStr := c.DefaultQuery("price_range", "")
	searchResult, err := db.SearchBook(title, authorName, priceRangeStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, searchResult)
}

func GetStoreInfo(c *gin.Context) {
	updateUserData(c, "GET")
	db, err := dal.NewBookStorer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	count, diffAuthors, err := db.GetStoreInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "There are "+
		strconv.FormatInt(count, 10)+" books by "+fmt.Sprintf("%.0f", *diffAuthors)+" different authors")

}

func GetUserData(c *gin.Context) {
	var r dal.Redis
	r.InitClient()
	username := c.DefaultQuery("username", "")
	userData := r.GetUserActivity(username)
	if userData ==nil{
		return
	}
	userPath0 := fmt.Sprintf("%v", userData[0])
	userPath1 := fmt.Sprintf("%v", userData[1])
	userPath2 := fmt.Sprintf("%v", userData[2])
	c.JSON(http.StatusOK,userPath0 + " , "+userPath1 + " , "+userPath2)
}

func updateUserData(c *gin.Context, op string) {
	username := c.DefaultQuery("username", "")
	r := dal.NewUserDater()
	req := c.Request.RequestURI
	out := fmt.Sprint(op, " ", req)
	r.AddActivity(username, out)
}
