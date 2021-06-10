package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/YuvalSteiy/book_service/models"
	"github.com/olivere/elastic/v7"
)

const INDEX_NAME = "books_yuval"
const PORT = 9200
const DOMAIN_URL = "http://es-search-7.fiverrdev.com"

type Elastic struct {
	Client *elastic.Client
}

func (e *Elastic) InitClient() error {
	setURL := fmt.Sprint(DOMAIN_URL, ":", PORT)
	client, err := elastic.NewClient(elastic.SetURL(setURL))
	if err != nil {
		return err
	}
	e.Client = client
	return nil
}

func (e *Elastic) GetBookByID(id string) (*models.Book, error) {
	ctx := context.Background()
	book, err := (e.Client).Get().Index(INDEX_NAME).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}
	var parsedBook models.Book
	json.Unmarshal(book.Source, &parsedBook)
	return &parsedBook, nil
}

func (e *Elastic) InsertBook(book *models.Book) (string, error) {
	ctx := context.Background()
	putStatus, err := (e.Client).Index().Index(INDEX_NAME).BodyJson(*book).Do(ctx)
	if err != nil {
		return "", err
	}
	return putStatus.Id, nil
}
func (e *Elastic) UpdateBook(title string, id string) error {
	ctx := context.Background()
	_, err := (e.Client).Update().Index(INDEX_NAME).Id(id).Doc(map[string]interface{}{"title": title}).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (e *Elastic) DeleteBook(id string) error {
	ctx := context.Background()
	_, err := (e.Client).Delete().Index("books_yuval").Id(id).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *Elastic) SearchBook(title string, authorName string, priceRangeStr string) ([]models.Book, error) {
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	if FieldExist(title) {
		matchQuery1 := elastic.NewMatchQuery("title", title)
		query.Filter(matchQuery1)
	}
	if FieldExist(authorName) {
		matchQuery2 := elastic.NewMatchQuery("author_name", authorName)
		query.Filter(matchQuery2)
	}
	if FieldExist(priceRangeStr) {
		priceRange := GetPriceRange(priceRangeStr)
		if priceRange != nil {
			rangeQuery := elastic.NewRangeQuery("price").From(priceRange[0]).To(priceRange[1])
			query.Filter(rangeQuery)
		}
	}
	searchResult, err := (e.Client).Search().Index(INDEX_NAME).Query(query).Size(100).Do(ctx)
	if err != nil {
		return nil, err
	}
	numHits := searchResult.Hits.TotalHits.Value
	resultsArr := make([]models.Book, numHits)
	for i, hit := range searchResult.Hits.Hits {
		json.Unmarshal(hit.Source, &resultsArr[i])
	}
	return resultsArr, nil
}

func (e *Elastic) GetStoreInfo() (int64, *float64, error) {
	ctx := context.Background()
	countService := elastic.NewCountService(e.Client)
	countResult, err := countService.Index(INDEX_NAME).Do(ctx)
	if err != nil {
		return 0, nil, err
	}
	agg := elastic.NewCardinalityAggregation().Field("author_name.keyword")
	diffAuthors, err := (e.Client.Search()).Index(INDEX_NAME).Aggregation("diff_authors", agg).Do(ctx)
	if err != nil {
		return 0, nil, err
	}
	numDiffAuthors, _ := diffAuthors.Aggregations.Cardinality("diff_authors")
	return countResult, numDiffAuthors.Value, nil
}

