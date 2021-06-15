package service

import (
	"github.com/YuvalSteiy/book_service/data_store"
	"github.com/pkg/errors"
)

func InitDataStore(){
	db := data_store.NewBookStorer()
	if db == nil{
		panic(errors.New("Could Not Init BookStore Database"))
	}
	userData := data_store.NewUserDater()
	if userData == nil {
		panic(errors.New("Could Not Init User Data Cache"))
	}
}