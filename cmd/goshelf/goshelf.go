/*
Goshelf is a radical book and book collection storage software!

Usage: goshelf [FLAGS] [CLI Command]
        -a Run in API mode, default false
        -dh Database address, default 0.0.0.0
        -dn Database name, default postgres
        -dp Database port, default 5432
        -dpw Database password, default ""
        -ds Database SSL mode, default postgres
        -du Database user, default postgres
        -p API mode: Host port, default 8080
        -s API mode: Host address, default 0.0.0.0

CLI Commands
        BOOKCREATE
        BOOKGET
        BOOKREMOVE
        BOOKFILTER
        COLLECTIONCREATE
        COLLECTIONGET
        COLLECTIONREMOVE

API Usage (flag -a)
        See API documentation for further help
*/

package goshelf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Max-Clark/goshelf/cmd/db"
	pg "github.com/Max-Clark/goshelf/cmd/db/postgresql"
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
)

const SchemaVersion = "v1"

type GoshelfConfig struct {
	RunApi   bool
	Host     string
	Port     int
	DbConfig db.ConnectionConfig
	Goshelf  GoshelfQuerier
}

type GoshelfQuerier interface {
	Connect() error
	BookCreate(b *v1.Book) (*int, error)
	BookGet(id int) (*v1.Book, error)
	BookRemove(id int) error
	BookFilter(title *string, genre *string, edition *int) ([]v1.Book, error)
	CollectionCreate(title *string, bookIds []int) (*string, error)
	CollectionGet(title *string) (*v1.Collection, error)
	CollectionRemove(title *string) error
}

func ApiStart(cfg GoshelfConfig) {
	StartServer(cfg)
}

func CliStart(cfg GoshelfConfig, flagSet *flag.FlagSet) {
	noFlagArgs := flagSet.Args()
	fMap := GetCliFuncMap()

	for _, v := range noFlagArgs {
		f, ok := fMap[strings.ToUpper(v)]

		// If the function given exists, run it
		if ok {
			f(&cfg)
			return
		}
	}

	w := os.Stderr
	fmt.Fprintln(w, "Invalid, missing, or unrecognized CLI command")
	PrintUsage(w, flagSet)
}

func Goshelf(args []string) {
	cfg, flagSet, err := InitFlags(os.Args)

	flagSet.Parse(args)

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
