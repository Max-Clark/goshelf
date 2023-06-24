package main

import (
	"log"
	"os"

	"github.com/Max-Clark/goshelf/cmd/db"
	"github.com/Max-Clark/goshelf/cmd/goshelf"
)

func main() {
	cfg, err := goshelf.InitFlags(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	cfg.Goshelf = db.PgDb{
		Config:        cfg.DbConfig,
		SchemaVersion: "v1", // TODO: make this dynamic & migrations
	}

	cfg.Goshelf.Connect()

	if cfg.RunApi {

	} else {

	}

	// TODO: Set up CLI/API
}
