package main

import (
	"fmt"
	"log"

	"github.com/Max-Clark/goshelf/cmd/cli"
)

func main() {
	str, err := cli.GetCliPrompt("Hi there: ")

	if err != nil {
		log.Fatal("error parsing flags", err)
	}

	fmt.Print(*str)

	// TODO: Set up CLI/API
}
