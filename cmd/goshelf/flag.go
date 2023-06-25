package goshelf

import (
	"flag"
	"fmt"
	"io"
)

func PrintUsage(w io.Writer, flagSet *flag.FlagSet) {
	fmt.Fprintln(w, "Usage: goshelf [FLAGS] [CLI Command]")

	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(w, "\t-%s %v", f.Name, f.Usage)
		fmt.Fprintln(w)
	})

	fmt.Fprintln(w)
	fmt.Fprintln(w, "CLI Commands")

	for cmd := range GetCliFuncMap() {
		fmt.Fprintln(w, "\t"+cmd)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "API Usage (flag -a)")
	fmt.Fprintln(w, "\tSee API documentation for further help")
}

func GetFlagSet(args []string, cfg *GoshelfConfig) *flag.FlagSet {
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

	return gsFlagSet
}

func InitFlags(args []string) (*GoshelfConfig, *flag.FlagSet, error) {
	cfg := GoshelfConfig{}
	gsFlagSet := GetFlagSet(args, &cfg)

	var err error
	if len(args) < 1 {
		err = gsFlagSet.Parse([]string{})
	} else {
		err = gsFlagSet.Parse(args[1:])
	}

	if err != nil {
		return nil, nil, err
	}

	return &cfg, gsFlagSet, nil
}
