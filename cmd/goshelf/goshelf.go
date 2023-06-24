package goshelf

import (
	"flag"
	"fmt"
	"log"
	"os"

	pg "github.com/Max-Clark/goshelf/cmd/db/postgresql"
	"github.com/Max-Clark/goshelf/cmd/http"
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
)

type GoshelfQuerier interface {
	Connect() error
	BookCreate(*v1.Book) (*int, error)
	BookGet(int) (*v1.Book, error)
	BookRemove(int) error
	BookFilter(*string, *string, *int) ([]v1.Book, error)
	CollectionCreate(*string, []int) (*string, error)
	CollectionGet(*string) (*v1.Collection, error)
	CollectionRemove(string) error
}

func ApiStart(cfg GoshelfConfig) error {
	http.StartServer(cfg.Host+":"+fmt.Sprint(cfg.Port), []http.PathFunction{})
	return nil
}

func CliStart(cfg GoshelfConfig, flagSet *flag.FlagSet) {
	noFlagArgs := flagSet.Args()
	fMap := GetCliFuncMap()

	for _, v := range noFlagArgs {
		f, ok := fMap[v]

		// If the function given exists, run it
		if ok {
			f(&cfg)
			return
		}
	}

	w := os.Stderr
	fmt.Fprintln(w, "Invalid, missing, or unrecognized CLI command")
	PrintFlagUsage(w, flagSet)
}

func Goshelf(args []string) {
	cfg, flagSet, err := InitFlags(os.Args)

	flagSet.Parse(args)

	fmt.Println(cfg)

	if err != nil {
		log.Fatal(err)
	}

	cfg.Goshelf = &pg.PgDb{
		Config:        cfg.DbConfig,
		SchemaVersion: "v1", // TODO: make this dynamic & migrations
	}

	cfg.Goshelf.Connect()

	if cfg.RunApi {
		ApiStart(*cfg)
	} else {
		CliStart(*cfg, flagSet)
	}
}
