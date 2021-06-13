package service

import (
	"fmt"
	"github.com/YuvalSteiy/book_service/data_store"
	"github.com/YuvalSteiy/book_service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	INIT_BOOKSTORE_FAIL = "Book Store Failed To Initialize"
	INIT_USERDATA_FAIL  = "User Data Failed To Initialize"
)

func GetBookByID(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_BOOKSTORE_FAIL))
		return
	}

	id := c.Param("id")
	book, err := db.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	c.JSON(http.StatusOK, book)
	return
}

func InsertBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_BOOKSTORE_FAIL))
		return
	}

	var book models.Book
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	id, err := db.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": fmt.Sprintf("Book inserted at ID: %s", id),
	})
}

func UpdateBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_BOOKSTORE_FAIL))
		return
	}

	id := c.Param("id")
	var book models.Book
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	err = db.UpdateBook(book.Title, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": fmt.Sprintf("Book updated at ID: %d", id),
	})
}

func DeleteBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_BOOKSTORE_FAIL))
		return
	}

	id := c.Param("id")
	err := db.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": fmt.Sprintf("Book deleted at ID: %d", id),
	})
}

func SearchBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_BOOKSTORE_FAIL))
		return
	}

	title := c.DefaultQuery("title", "")
	authorName := c.DefaultQuery("author_name", "")
	priceRangeStr := c.DefaultQuery("price_range", "")

	searchResult, err := db.SearchBook(title, authorName, priceRangeStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	c.JSON(http.StatusOK, searchResult)
}

func GetStoreInfo(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_BOOKSTORE_FAIL))
		return
	}

	count, diffAuthors, err := db.GetStoreInfo()
	if diffAuthors == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat("Failed To Retrieve Store Info"))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}

	response := fmt.Sprintf("There are %d books by %.0f different authors", count, *diffAuthors)
	c.JSON(http.StatusOK, map[string]string{
		"response": response,
	})

}

func GetUserData(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	if username == "" {
		return
	}

	r := data_store.NewUserDater()
	if r == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_USERDATA_FAIL))
	}

	response, err := r.GetUserActivity(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(err.Error()))
		return
	}
	if response == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat("Got Nil response"))
		return
	}

	userPath1 := fmt.Sprintf("%v", response[0])
	userPath2 := fmt.Sprintf("%v", response[1])
	userPath3 := fmt.Sprintf("%v", response[2])

	userData := map[string]string{"last action": userPath1}
	if userPath2 != "" {
		userData["2nd last action"] = userPath2
		if userPath3 != "" {
			userData["3rd last action"] = userPath3
		}
	}

	c.JSON(http.StatusOK, userData)
}

func updateUserData(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	if username == "" {
		return
	}

	r := data_store.NewUserDater()
	if r == nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat(INIT_USERDATA_FAIL))
	}

	req := c.Request.RequestURI
	method := c.Request.Method
	activity := fmt.Sprintf("%s %s", method, req)
	err := r.AddActivity(username, activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputErrorFormat("Failed To Update User Data"))
	}
}

func outputErrorFormat(err string) map[string]string {
	return map[string]string{
		"Error": err,
	}
}
