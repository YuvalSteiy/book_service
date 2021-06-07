package bl

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"test/dal"
)

type db interface{
	DBInitClient() error
	DBGetBookByID(id string) (*elastic.GetResult,error)
	DBGetStoreData() (int64, *elastic.SearchResult,error)
	DBPutBook(book *dal.Book) (*elastic.IndexResponse,error)
	DBUpdateBook(title string, id string) (*elastic.UpdateResponse,error)
	DBDeleteBook(id string) (*elastic.DeleteResponse,error)
	DBSearchBook(title string, authorName string, priceRange [2]float64, opcode int) (*elastic.SearchResult,error)
}
type userDataCache interface{
	InitClient()
	CacheAddActivity(username string, req string)error
	CacheCreateNewClientActivity(username string, req string,client *redis.Client)
	CacheUpdateClientActivity(username string, req string, client *redis.Client)
	CacheGetUserActivity(username string) ([]interface{},error)
}

func GetBookByID(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	updateUserData(username,"GET",c.Request.RequestURI)
	id := c.Param("id")
	db := dal.Elastic{}
	err := db.DBInitClient()
	if err != nil{
		handleError(err,c)
	}
	book, err := db.DBGetBookByID(id)
	if err!=nil{
		handleError(err,c)
	}
	c.JSON(http.StatusOK, book.Source)
}

func PutBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	updateUserData(username,"PUT",c.Request.RequestURI)
	jsonData,err := ioutil.ReadAll(c.Request.Body)
	if err!=nil{
		handleError(err,c)
	}
	book := dal.Book{}
	json.Unmarshal(jsonData,&book)
	db := dal.Elastic{}
	err = db.DBInitClient()
	if err != nil{
		handleError(err,c)
	}
	putStatus,err := db.DBPutBook(&book)
	if err!=nil{
		handleError(err,c)
	}
	c.JSON(http.StatusOK, putStatus.Id)
}

func UpdateBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	updateUserData(username,"POST",c.Request.RequestURI)
	jsonData,err := ioutil.ReadAll(c.Request.Body)
	if err!=nil{
		handleError(err,c)	}
	book := dal.Book{}
	json.Unmarshal(jsonData,&book)
	db := dal.Elastic{}
	err = db.DBInitClient()
	if err != nil{
		handleError(err,c)
	}
	response,err := db.DBUpdateBook(book.Title,c.Param("id"))
	if err!=nil{
		handleError(err,c)
	}
	c.JSON(http.StatusOK, response.Id)
}

func DeleteBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	updateUserData(username,"Delete",c.Request.RequestURI)
	id := c.Param("id")
	db := dal.Elastic{}
	err := db.DBInitClient()
	if err != nil{
		handleError(err,c)
	}
	response,err := db.DBDeleteBook(id)
	if err!=nil{
		handleError(err,c)
	}
	c.JSON(http.StatusOK,response)
}

func SearchBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	updateUserData(username,"GET",c.Request.RequestURI)
	title       := c.DefaultQuery("title","no_value")
	authorName  := c.DefaultQuery("author_name","no_value")
	priceRangeS := c.DefaultQuery("price_range", "no_value")
	db  := dal.Elastic{}
	err := db.DBInitClient()
	if err != nil{
		handleError(err,c)
	}
	priceRange := getPriceRange(priceRangeS, c)
	var opcode int = 0
	if fieldExist(priceRangeS){
		opcode +=1
	}
	if fieldExist(title){
		opcode +=2
	}
	if fieldExist(authorName){
		opcode +=4
	}
	searchResult,err := db.DBSearchBook(title, authorName, priceRange, opcode)
	if err!=nil{
		handleError(err,c)
	}
	for _, hit := range searchResult.Hits.Hits {
		c.JSON(http.StatusOK, hit.Source)
	}
}

func GetStoreData(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	updateUserData(username,"GET",c.Request.RequestURI)
	db  := dal.Elastic{}
	err := db.DBInitClient()
	if err != nil{
		handleError(err,c)
	}
	count,diffAuthors,err := db.DBGetStoreData()
	if err!=nil{
		handleError(err,c)
	}
	c.JSON(http.StatusOK, count)
	c.JSON(http.StatusOK, diffAuthors.Hits.TotalHits)
}
func GetUserData(c *gin.Context){
	r := dal.Redis{}
	r.InitClient()
	username    := c.DefaultQuery("username","no_value")
	if !fieldExist(username){
		return
	}
	userData,err := r.CacheGetUserActivity(username)
	if err!=nil{
		handleError(err,c)
	}
	c.JSON(http.StatusOK, userData[0])
	c.JSON(http.StatusOK," , ")
	if userData[1] != "" {
		c.JSON(http.StatusOK, userData[1])
	}
	if userData[2] != "" {
		c.JSON(http.StatusOK, userData[2])
	}


}

func fieldExist(check string) (bool){
	if check == "no_value"{
		return false
	}
	return true
}

func getPriceRange(priceRangeS string, c *gin.Context) ([2]float64){
	priceRange := [2]float64{0.0,0.0}
	var err error
	if !fieldExist(priceRangeS){
		return priceRange
	}
	priceRange[0],err = strconv.ParseFloat(strings.Split(priceRangeS,"-")[0], 64)
	if err!=nil{
		handleError(err,c)
	}
	priceRange[1],err = strconv.ParseFloat(strings.Split(priceRangeS,"-")[1], 64)
	if err!=nil{
		handleError(err,c)
	}
	return priceRange
}

func updateUserData(username string, op string, req string){
	r := dal.Redis{}
	r.InitClient()
	if !fieldExist(username){
		return
	}
	out := fmt.Sprint(op," ",req)
	r.CacheAddActivity(username, out)
}

func handleError(err error, c *gin.Context){
	fmt.Println(err.Error())
	c.JSON(http.StatusInternalServerError,err)
}
