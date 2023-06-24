package goshelf

import (
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
)

type GoshelfQuerier interface {
	Connect() error
	BookCreate(*v1.Book) (*int, error)
	BookGet(int) (*v1.Book, error)
	BookRemove(int) error
	BookFilter(*v1.Book) ([]v1.Book, error)
	CollectionCreate(*[]int) (*int, error)
	CollectionGet(int) (*v1.Collection, error)
	CollectionRemove(int) error
}

func ApiStart(cfg GoshelfConfig) {
	// http.StartServer(cfg.Host + ":" + fmt.Sprint(cfg.Port))
}
