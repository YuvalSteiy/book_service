package data_store

import (
	"github.com/YuvalSteiy/book_service/models"
<<<<<<< HEAD
	"sync"
)

var bookStorer BookStorer
var bookStorerOnce sync.Once
=======
)

var bookStorer BookStorer
>>>>>>> 3307927fd33aca1aeb0a88015e2937965045df37

type BookStorer interface {
	GetBookByID(id string) (*models.Book, error)
	InsertBook(book models.Book) (string, error)
	UpdateBook(title string, id string) error
	DeleteBook(id string) error
	SearchBook(title string, authorName string, priceRangeStr string) ([]models.Book, error)
	GetStoreInfo() (int64, *float64, error)
}

func NewBookStorer() BookStorer {
<<<<<<< HEAD
	bookStorerOnce.Do(func(){
		bookStorer = CreateElastic()
	})
=======
	if bookStorer == nil {
		bookStorer = CreateElastic()
	}
>>>>>>> 3307927fd33aca1aeb0a88015e2937965045df37

	return bookStorer
}
