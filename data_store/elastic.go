package data_store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/YuvalSteiy/book_service/models"
	"github.com/olivere/elastic/v7"
)

const (
	INDEX_NAME = "books_yuval"
	PORT       = 9200
	DOMAIN_URL = "http://es-search-7.fiverrdev.com"
)

type Elastic struct {
	Client *elastic.Client
}

func initElasticClient() *elastic.Client {
	setURL := fmt.Sprintf("%s:%d", DOMAIN_URL, PORT)
	client, err := elastic.NewClient(elastic.SetURL(setURL))
	if err != nil {
		return nil
	}

	return client
}

func CreateElastic() *Elastic {
	e := Elastic{Client: initElasticClient()}
	if e.Client == nil {
		return nil
	}

	return &e
}

func (e *Elastic) GetBookByID(id string) (*models.Book, error) {
	ctx := context.Background()
	response, err := (e.Client).Get().Index(INDEX_NAME).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	var book models.Book
	err = json.Unmarshal(response.Source, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (e *Elastic) InsertBook(book models.Book) (string, error) {
	ctx := context.Background()
	response, err := (e.Client).Index().Index(INDEX_NAME).BodyJson(book).Do(ctx)
	if err != nil {
		return "", err
	}

	return response.Id, nil
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
	_, err := (e.Client).Delete().Index(INDEX_NAME).Id(id).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e *Elastic) SearchBook(title string, authorName string, priceRangeStr string) ([]models.Book, error) {
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	if title != "" {
		query.Filter(elastic.NewMatchQuery("title", title))
	}
	if authorName != "" {
		query.Filter(elastic.NewMatchQuery("author_name", authorName))
	}
	if priceRangeStr != "" {
		priceRange, err := GetPriceRange(priceRangeStr)
		if err != nil {
			return nil, err
		}
		if priceRange == nil {
			return nil, nil
		}
		query.Filter(elastic.NewRangeQuery("price").From(priceRange[0]).To(priceRange[1]))
	}
	searchResult, err := (e.Client).Search().Index(INDEX_NAME).Query(query).Size(100).Do(ctx)
	if err != nil {
		return nil, err
	}

	numHits := searchResult.Hits.TotalHits.Value
	resultsArr := make([]models.Book, numHits)
	for i, hit := range searchResult.Hits.Hits {
		err = json.Unmarshal(hit.Source, &resultsArr[i])
		if err != nil {
			return nil, err
		}
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
	response, err := (e.Client.Search()).Index(INDEX_NAME).Aggregation("diff_authors", agg).Do(ctx)
	if err != nil {
		return 0, nil, err
	}

	numDiffAuthors, found := response.Aggregations.Cardinality("diff_authors")
	if found == false {
		return 0, nil, nil
	}

	return countResult, numDiffAuthors.Value, nil
}
