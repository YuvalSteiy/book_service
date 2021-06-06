package bl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"test/dal"
)

func GetBookByID(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	id := c.Param("id")
	client := dal.ElasticInitClient()
	get1 := dal.ElasticGetBookByID(client, id)
	updateUserData(username,"GET " +c.Request.RequestURI)
	c.JSON(http.StatusOK, get1.Source)
}

func PutBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	jsonData,err := ioutil.ReadAll(c.Request.Body)
	if err!=nil{
		panic(err)
	}
	book := dal.Book{}
	json.Unmarshal(jsonData,&book)
	client := dal.ElasticInitClient()
	put1 := dal.ElasticPutBook(client, &book)
	updateUserData(username,"PUT " +c.Request.RequestURI)
	c.JSON(http.StatusOK, put1.Id)
}

func UpdateBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	jsonData,err := ioutil.ReadAll(c.Request.Body)
	if err!=nil{
		panic(err)
	}
	book := dal.Book{}
	json.Unmarshal(jsonData,&book)
	client := dal.ElasticInitClient()
	response := dal.ElasticUpdateBook(client,book.Title,c.Param("id"))
	updateUserData(username,"POST " +c.Request.RequestURI)
	c.JSON(http.StatusOK, response.Id)
}

func DeleteBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	id := c.Param("id")
	client := dal.ElasticInitClient()
	response := dal.ElasticDeleteBook(client, id)
	updateUserData(username,"Delete " +c.Request.RequestURI)
	c.JSON(http.StatusOK,response)
}

func SearchBook(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	title       := c.DefaultQuery("title","no_value")
	authorName  := c.DefaultQuery("author_name","no_value")
	priceRangeS := c.DefaultQuery("price_range", "no_value")
	client      := dal.ElasticInitClient()
	price_range := getPriceRange(priceRangeS)
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
	searchResult := dal.ElasticSearchBook(client, title, authorName, price_range, opcode)
	updateUserData(username,"GET " +c.Request.RequestURI)
	// Iterate through results
	for _, hit := range searchResult.Hits.Hits {
		c.JSON(http.StatusOK, hit.Source)
	}
}

func GetStoreData(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	client := dal.ElasticInitClient()
	count,diffAuthors := dal.ElasticGetStoreData(client)
	updateUserData(username,"GET " +c.Request.RequestURI)
	c.JSON(http.StatusOK, count)
	c.JSON(http.StatusOK, diffAuthors.Hits.TotalHits)
}
func GetUserData(c *gin.Context){
	username    := c.DefaultQuery("username","no_value")
	if username == "no_value"{
		panic("No Username Given")
	}
	userData := dal.ReddisGetUserRecord(username)
	c.JSON(http.StatusOK, userData)
}

func fieldExist(check string) (bool){
	if check == "no_value"{
		return false
	}
	return true
}

func getPriceRange(priceRangeS string) ([2]float64){
	priceRange := [2]float64{0.0,0.0}
	var err error
	if !fieldExist(priceRangeS){
		return priceRange
	}
	priceRange[0],err = strconv.ParseFloat(strings.Split(priceRangeS,"-")[0], 64)
	if err!=nil{
		panic(err)
	}
	priceRange[1],err = strconv.ParseFloat(strings.Split(priceRangeS,"-")[1], 64)
	if err!=nil{
		panic(err)
	}
	return priceRange
}

func updateUserData(username string, req string){
	if fieldExist(username){
		dal.ReddisAddRecord(username, req)
	}
}
