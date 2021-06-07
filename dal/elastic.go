package dal

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)
const INDEXNAME = "books_yuval"
const PORT = 9200
const DOMAINURL = "http://es-search-7.fiverrdev.com"
type Book struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	EbookAvailable bool    `json:"ebook_available"`
	Price          float64 `json:"price"`
	PublishDate    string  `json:"publish_date"`
}

type Elastic struct{
	Client *elastic.Client
}

func (e *Elastic) DBInitClient() error {
	setURL := fmt.Sprint(DOMAINURL,":",PORT)
	client, err := elastic.NewClient(elastic.SetURL(setURL))
	if err != nil {
		return err
	}
	e.Client = client
	return nil
}

func (e *Elastic) DBGetBookByID(id string) (*elastic.GetResult,error){
	ctx := context.Background()
	book, err := (e.Client).Get().Index(INDEXNAME).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (e *Elastic) DBGetStoreData() (int64, *elastic.SearchResult,error){
	ctx := context.Background()
	//get document count in index
	countService := elastic.NewCountService(e.Client)
	countResult, err := countService.Index(INDEXNAME).Do(ctx)
	if err != nil {
		return 0, nil, err
	}
	//get number of distinct authors
	agg := elastic.NewCardinalityAggregation().Field("author_name.keyword")
	diffAuthors, err := (e.Client.Search()).Index(INDEXNAME).Aggregation("diff_authors", agg).Do(ctx)
	if err !=nil{
		return 0,nil,err
	}
	return countResult,diffAuthors,nil

}

func (e *Elastic) DBPutBook(book *Book) (*elastic.IndexResponse,error){
	ctx := context.Background()
	// Index a tweet (using JSON serialization)
	putStatus, err := (e.Client).Index().Index(INDEXNAME).BodyJson(*book).Do(ctx)
	if err!=nil{
		return nil, err
	}
	return putStatus,nil
}

func (e *Elastic) DBUpdateBook(title string, id string) (*elastic.UpdateResponse,error){
	ctx := context.Background()
	response, err := (e.Client).Update().Index(INDEXNAME).Id(id).Doc(map[string]interface{}{"title": title}).Do(ctx)
	if err !=nil{
		return nil, err
	}
	return response,nil
}

func (e *Elastic) DBDeleteBook(id string) (*elastic.DeleteResponse,error){
	ctx := context.Background()
	response, err := (e.Client).Delete().Index("books_yuval").Id(id).Do(ctx)
	if err != nil{
		return nil, err
	}
	return response,nil
}

func (e *Elastic) DBSearchBook(title string, authorName string, priceRange [2]float64, opcode int) (*elastic.SearchResult,error){
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	matchQuery1 := elastic.NewMatchQuery("title", title)
	matchQuery2 := elastic.NewMatchQuery("author_name", authorName)
	rangeQuery := elastic.NewRangeQuery("price").From(priceRange[0]).To(priceRange[1])
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
	searchResult, err := (e.Client).Search().Index(INDEXNAME).Query(searchQuery).Size(100).Do(ctx)
	if err != nil{
		return nil, err
	}
	return searchResult,nil
}

