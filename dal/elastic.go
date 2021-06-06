package dal

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
)
const INDEXNAME = "books_yuval"
// Book struct for deserialization of documents from json
type Book struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	EbookAvailable bool    `json:"ebook_available"`
	Price          float64 `json:"price"`
	PublishDate    string  `json:"publish_date"`
}

func ElasticInitClient() *elastic.Client {
	// Pass connection port
	port := 9200
	strPort := strconv.Itoa(port)
	// Pass domain string
	domainURL := "http://es-search-7.fiverrdev.com"
	// Complete URL
	setURL := domainURL + ":" + strPort
	// Instantiate a client (elastic library instance)
	client, err := elastic.NewClient(elastic.SetURL(setURL))
	if err != nil {
		fmt.Println("elastic.NewClient() ERROR:", err)
		log.Fatalf("quitting connection..")
	}
	return client
}

func ElasticGetBookByID(client *elastic.Client, id string) *elastic.GetResult{
	ctx := context.Background()
	get1, err := client.Get().Index(INDEXNAME).Id(id).Do(ctx)
	if err != nil {
		panic(err)
	}
	return get1
}

func ElasticGetStoreData(client *elastic.Client) (int64, *elastic.SearchResult){
	ctx := context.Background()
	//get document count in index
	countService := elastic.NewCountService(client)
	countResult, err := countService.Index(INDEXNAME).Do(ctx)
	if err != nil {
		panic(err)
	}

	//get number of distinct authors
	agg := elastic.NewCardinalityAggregation().Field("author_name.keyword")
	diffAuthors, err := client.Search().Index(INDEXNAME).Aggregation("diff_authors", agg).Do(ctx)
	if err !=nil{
		panic(err)
	}
	return countResult,diffAuthors

}

func ElasticPutBook(client *elastic.Client, book *Book) *elastic.IndexResponse{
	ctx := context.Background()
	// Index a tweet (using JSON serialization)
	put1, err := client.Index().Index(INDEXNAME).BodyJson(*book).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	return put1
}

func ElasticUpdateBook(client *elastic.Client, title string, id string) *elastic.UpdateResponse{
	ctx := context.Background()
	response, err := client.Update().Index(INDEXNAME).Id(id).Doc(map[string]interface{}{"title": title}).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	return response
}

func ElasticDeleteBook(client *elastic.Client, id string) *elastic.DeleteResponse{
	ctx := context.Background()
	response, err := client.Delete().Index("books_yuval").Id(id).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	return response
}

func ElasticSearchBook(client *elastic.Client, title string, author_name string, price_range [2]float64, opcode int) *elastic.SearchResult{
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	matchQuery1 := elastic.NewMatchQuery("title", title)
	matchQuery2 := elastic.NewMatchQuery("author_name", author_name)
	rangeQuery := elastic.NewRangeQuery("price").From(price_range[0]).To(price_range[1])
	var searchQuery *elastic.BoolQuery
	switch opcode{//opcode specifies which arguments were given to query
		case 0://no search arguments
			searchQuery =query.Must()
		case 1://only price range specified
			searchQuery = query.Must(rangeQuery)
		case 2://only title
			searchQuery = query.Must(matchQuery1)
		case 3://price and title
			searchQuery = query.Must(matchQuery1, rangeQuery)
		case 4://only author
			searchQuery = query.Must(matchQuery2)
		case 5://price and author
			searchQuery = query.Must(matchQuery2, rangeQuery)
		case 6://title and author
			searchQuery = query.Must(matchQuery1, matchQuery2)
		case 7://title author and price
			searchQuery = query.Must(matchQuery1, matchQuery2, rangeQuery)
	}
	searchResult, err := client.Search().Index(INDEXNAME).Query(searchQuery).Size(100).Do(ctx)
	if err != nil {
		panic(err)
	}
		return searchResult
}

