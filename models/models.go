package models

type Book struct {
	Title          *string  `json:"title"`
	AuthorName     *string  `json:"author_name"`
	EbookAvailable *bool    `json:"ebook_available"`
	Price          *float64 `json:"price"`
	PublishDate    *string  `json:"publish_date"`
}
