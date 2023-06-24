package goshelf

import (
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
)

type GoshelfQuerier interface {
	Connect() error
	BookCreate(v1.Book) error
	BookGet(uint32) (*v1.Book, error)
	BookRemove(uint32) error
	BookFilter(v1.Book) ([]v1.Book, error)
	CollectionCreate([]uint32) error
	CollectionGet(uint32) (*v1.Collection, error)
	CollectionRemove(v1.Collection) error
}

func ApiStart(cfg GoshelfConfig) {
	// http.StartServer(cfg.Host + ":" + fmt.Sprint(cfg.Port))
}
