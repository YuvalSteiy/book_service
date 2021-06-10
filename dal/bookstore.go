package dal

import (
	"github.com/YuvalSteiy/book_service/models"
)

type BookStorer interface {
	InitClient() error
	GetBookByID(id string) (*models.Book, error)
	InsertBook(book *models.Book) (string, error)
	UpdateBook(title string, id string) (error)
	DeleteBook(id string) (error)
	SearchBook(title string, authorName string, priceRangeStr string) ([]models.Book, error)
	GetStoreInfo() (int64, *float64, error)
}

func NewBookStorer() (BookStorer, error) {
	var db Elastic
	err := db.InitClient()
	return &db, err
}
