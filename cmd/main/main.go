package main

import (
	"fmt"
	"log"
	"os"

	pg "github.com/Max-Clark/goshelf/cmd/db/postgresql"
	"github.com/Max-Clark/goshelf/cmd/goshelf"
	"github.com/Max-Clark/goshelf/cmd/http"
)

func main() {
	cfg, err := goshelf.InitFlags(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	cfg.Goshelf = &pg.PgDb{
		Config:        cfg.DbConfig,
		SchemaVersion: "v1", // TODO: make this dynamic & migrations
	}

	cfg.Goshelf.Connect()

	if cfg.RunApi {
		http.StartServer(cfg.Host+":"+fmt.Sprint(cfg.Port), []http.PathFunction{})
	} else {

	}

	// TODO: Set up CLI/API
}
