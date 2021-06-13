package service

import (
	"fmt"
	"github.com/YuvalSteiy/book_service/data_store"
	"github.com/YuvalSteiy/book_service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	INIT_BOOKSTORE_FAIL = "Book Store failed to initialize"
	INIT_USERDATA_FAIL  = "User data failed to initialize"
	STORE_INFO_FAIL     = "Failed to retrieve store info"
	INALID_INPUT_FAIL   = "Invalid input"
)

func GetBookByID(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, INIT_BOOKSTORE_FAIL)
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
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, INIT_BOOKSTORE_FAIL)
		return
	}

	var book models.Book
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	id, err := db.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": "Book inserted at ID: " + id,
	})
}

func UpdateBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, INIT_BOOKSTORE_FAIL)
		return
	}

	id := c.Param("id")
	var book models.Book
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = db.UpdateBook(book.Title, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": "Book updated at ID: " + id,
	})
}

func DeleteBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, INIT_BOOKSTORE_FAIL)
		return
	}

	id := c.Param("id")
	err := db.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"response": "Book deleted at ID: " + id,
	})
}

func SearchBook(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, INIT_BOOKSTORE_FAIL)
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

	if searchResult == nil {
		c.JSON(http.StatusInternalServerError, INALID_INPUT_FAIL)
		return
	}

	c.JSON(http.StatusOK, searchResult)
}

func GetStoreInfo(c *gin.Context) {
	updateUserData(c)
	db := data_store.NewBookStorer()
	if db == nil {
		c.JSON(http.StatusInternalServerError, INIT_BOOKSTORE_FAIL)
		return
	}

	count, diffAuthors, err := db.GetStoreInfo()
	if diffAuthors == nil {
		c.JSON(http.StatusInternalServerError, STORE_INFO_FAIL)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
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
		c.JSON(http.StatusInternalServerError, INIT_USERDATA_FAIL)
	}

	response := r.GetUserActivity(username)
	if response == nil {
		return
	}

	userPath1 := fmt.Sprintf("%v", response[0])
	userPath2 := fmt.Sprintf("%v", response[1])
	userPath3 := fmt.Sprintf("%v", response[2])

	userData := map[string]string{"path1": userPath1}
	if userPath2 != "" {
		userData["path2"] = userPath2
		if userPath3 != "" {
			userData["path3"] = userPath3
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
		c.JSON(http.StatusInternalServerError, INIT_USERDATA_FAIL)
	}

	req := c.Request.RequestURI
	method := c.Request.Method
	activity := fmt.Sprintf("%s %s", method, req)
	r.AddActivity(username, activity)
}
