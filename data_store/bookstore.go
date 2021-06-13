package data_store

import (
	"github.com/YuvalSteiy/book_service/models"
	"sync"
)

var bookStorer BookStorer
var bookStorerOnce sync.Once

type BookStorer interface {
	GetBookByID(id string) (*models.Book, error)
	InsertBook(book models.Book) (string, error)
	UpdateBook(title string, id string) error
	DeleteBook(id string) error
	SearchBook(title string, authorName string, priceRangeStr string) ([]models.Book, error)
	GetStoreInfo() (int64, *float64, error)
}

func NewBookStorer() BookStorer {
	bookStorerOnce.Do(func() {
		bookStorer = CreateElastic()
	})

	return bookStorer
}
