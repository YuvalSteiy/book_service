package service

import (
	"github.com/YuvalSteiy/book_service/data_store"
	"github.com/pkg/errors"
)

func InitDataStore() error {
	db := data_store.NewBookStorer()
	if db == nil{
		return errors.New("Could Not Init BookStore Database")
	}
	userData := data_store.NewUserDater()
	if userData == nil {
		return errors.New("Could Not Init User Data Cache")
	}
	return nil
}