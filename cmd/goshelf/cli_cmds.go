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
	"BOOKCREATE":       CliBookCreate,
	"BOOKGET":          CliBookGet,
	"BOOKREMOVE":       CliBookRemove,
	"BOOKFILTER":       CliBookFilter,
	"COLLECTIONCREATE": CliCollectionCreate,
	"COLLECTIONGET":    CliCollectionGet,
	"COLLECTIONREMOVE": CliCollectionRemove,
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
	prompt := "\tEnter title: "
	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter author's first name: "
	aFirst, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter author's first name: "
	aLast, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter description (optional): "
	desc, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter edition (optional): "
	ed, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter genre (optional): "
	genre, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter publish date YYYY-MM-dd (optional): "
	date, err := cli.GetCliPrompt(&prompt, os.Stdin)
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
	prompt := "\tEnter book id: "
	idStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
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
	prompt := "\tEnter book id: "
	idStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	idInt64, err := strconv.ParseInt(*idStr, 10, 32)
	PanicErrorHandler(err)

	id := int(idInt64)

	err = cfg.Goshelf.BookRemove(id)
	PanicErrorHandler(err)
}

func CliBookFilter(cfg *GoshelfConfig) {
	prompt := "\tEnter partial title to search (optional): "
	titleStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter edition (optional): "
	edStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	prompt = "\tEnter partial genre to search (optional): "
	genreStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	var title *string
	var genre *string
	var edition *int

	title = nil
	genre = nil
	edition = nil

	if *titleStr != "" {
		title = titleStr
	}

	if *genreStr != "" {
		genre = genreStr
	}

	if *edStr != "" {
		edInt64, err := strconv.ParseInt(*edStr, 10, 32)
		PanicErrorHandler(err)

		editionInt := int(edInt64)
		edition = &editionInt
	}

	books, err := cfg.Goshelf.BookFilter(title, genre, edition)
	PanicErrorHandler(err)

	for _, book := range books {
		json, err := json.Marshal(book)
		PanicErrorHandler(err)

		fmt.Println(string(json))
	}
}

func CliCollectionCreate(cfg *GoshelfConfig) {
	prompt := "\tEnter collection title: "
	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	bookIds := []int{}

	for {
		prompt := "\tEnter book id (press enter when finished): "
		bookId, err := cli.GetIntFromCli(&prompt, os.Stdin, os.Stdout)
		PanicErrorHandler(err)

		if bookId == nil {
			break
		}

		bookIds = append(bookIds, *bookId)
	}

	_, err = cfg.Goshelf.CollectionCreate(title, bookIds)
	PanicErrorHandler(err)
}

func CliCollectionGet(cfg *GoshelfConfig) {
	prompt := "\tEnter collection title: "
	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	col, err := cfg.Goshelf.CollectionGet(title)
	PanicErrorHandler(err)

	json, err := json.Marshal(col)
	PanicErrorHandler(err)

	fmt.Println(string(json))
}

func CliCollectionRemove(cfg *GoshelfConfig) {
	prompt := "\tEnter collection title: "
	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
	PanicErrorHandler(err)

	err = cfg.Goshelf.CollectionRemove(title)
	PanicErrorHandler(err)
}
