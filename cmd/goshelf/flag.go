package goshelf

import (
	"flag"
	"fmt"
	"os"

	"github.com/Max-Clark/goshelf/cmd/db"
)

type GoshelfConfig struct {
	RunApi   bool
	Host     string
	Port     int
	DbConfig db.ConnectionConfig
	Goshelf  GoshelfQuerier
}

func PrintHelp() {
	fmt.Fprintf(os.Stderr, "Custom help %s:\n", os.Args[0])

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(os.Stderr, "    %v\n", f.Usage) // f.Name, f.Value
	})
}

func InitFlags(args []string) (*GoshelfConfig, error) {
	cfg := GoshelfConfig{}
	gsFlagSet := flag.NewFlagSet("gsFlagSet", flag.ContinueOnError)

	gsFlagSet.StringVar(&cfg.Host, "s", "0.0.0.0", "API mode: Host address, default 0.0.0.0")
	gsFlagSet.IntVar(&cfg.Port, "p", 8080, "API mode: Host port, default 8080")
	gsFlagSet.BoolVar(&cfg.RunApi, "a", false, "Run in API mode, default false")

	gsFlagSet.StringVar(&cfg.DbConfig.Host, "dh", "0.0.0.0", "Database address, default 0.0.0.0")
	gsFlagSet.IntVar(&cfg.DbConfig.Port, "dp", 5432, "Database port, default 5432")
	gsFlagSet.StringVar(&cfg.DbConfig.User, "du", "postgres", "Database user, default postgres")
	gsFlagSet.StringVar(&cfg.DbConfig.Password, "dpw", "", "Database password, default \"\"")
	gsFlagSet.StringVar(&cfg.DbConfig.DbName, "dn", "postgres", "Database name, default postgres")
	gsFlagSet.StringVar(&cfg.DbConfig.SslMode, "ds", "disable", "Database SSL mode, default postgres")

	var err error
	if len(args) < 1 {
		err = gsFlagSet.Parse([]string{})
	} else {
		err = gsFlagSet.Parse(args[1:])
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
