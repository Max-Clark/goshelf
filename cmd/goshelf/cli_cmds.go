package goshelf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Max-Clark/goshelf/cmd/cli"
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
)

var cliFuncMap = map[string]func(*GoshelfConfig){
	"bookcreate":       CliBookCreate,
	"bookget":          CliBookGet,
	"bookremove":       CliBookRemove,
	"bookfilter":       CliBookFilter,
	"collectioncreate": CliCollectionCreate,
	"collectionget":    CliCollectionGet,
	"collectionremove": CliCollectionRemove,
}

func GetCliFuncMap() map[string]func(*GoshelfConfig) {
	return cliFuncMap
}

func PanicErrorHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// TODO: automate these based on reflection
// Creates a book from the cli. Requires a fully configured Goshelfconfig.
func CliBookCreate(cfg *GoshelfConfig) {
	// TODO: implement cli tooling interface
	title, err := cli.GetCliPrompt("\tEnter title: ", os.Stdin)
	PanicErrorHandler(err)

	aFirst, err := cli.GetCliPrompt("\tEnter author's first name: ", os.Stdin)
	PanicErrorHandler(err)

	aLast, err := cli.GetCliPrompt("\tEnter author's last name: ", os.Stdin)
	PanicErrorHandler(err)

	desc, err := cli.GetCliPrompt("\tEnter description (optional): ", os.Stdin)
	PanicErrorHandler(err)

	ed, err := cli.GetCliPrompt("\tEnter edition (optional): ", os.Stdin)
	PanicErrorHandler(err)

	genre, err := cli.GetCliPrompt("\tEnter genre (optional): ", os.Stdin)
	PanicErrorHandler(err)

	date, err := cli.GetCliPrompt("\tEnter publish date YYYY-MM-dd (optional): ", os.Stdin)
	PanicErrorHandler(err)

	book := v1.Book{
		Title: *title,
		Author: v1.Author{
			FirstName: *aFirst,
			LastName:  *aLast,
		},
	}

	if *desc != "" {
		book.Description = desc
	}

	if *genre != "" {
		book.Genre = genre
	}

	if *ed != "" {
		edInt64, err := strconv.ParseInt(*ed, 10, 32)
		edInt := int(edInt64)

		if err != nil {
			log.Panic("invalid edition (must be integer)")
		}

		if *ed != "" {
			book.Edition = &edInt
		}
	}

	if *date != "" {
		pDate, err := time.Parse("2006-01-02", *date)

		if err != nil {
			log.Panic("invalid time format (must match YYYY-MM-dd)")
		}

		book.PublishDate = &pDate
	}

	id, err := cfg.Goshelf.BookCreate(&book)

	PanicErrorHandler(err)

	fmt.Printf("%v", *id)
}

func CliBookGet(cfg *GoshelfConfig) {
	idStr, err := cli.GetCliPrompt("\tEnter book id: ", os.Stdin)
	PanicErrorHandler(err)

	idInt64, err := strconv.ParseInt(*idStr, 10, 32)
	PanicErrorHandler(err)

	id := int(idInt64)

	book, err := cfg.Goshelf.BookGet(id)
	PanicErrorHandler(err)

	json, err := json.Marshal(book)
	PanicErrorHandler(err)

	fmt.Print(string(json))
}

func CliBookRemove(cfg *GoshelfConfig) {
	idStr, err := cli.GetCliPrompt("\tEnter book id: ", os.Stdin)
	PanicErrorHandler(err)

	idInt64, err := strconv.ParseInt(*idStr, 10, 32)
	PanicErrorHandler(err)

	id := int(idInt64)

	err = cfg.Goshelf.BookRemove(id)
	PanicErrorHandler(err)
}

func CliBookFilter(cfg *GoshelfConfig) {

}

func CliCollectionCreate(cfg *GoshelfConfig) {

}

func CliCollectionGet(cfg *GoshelfConfig) {

}

func CliCollectionRemove(cfg *GoshelfConfig) {

}
