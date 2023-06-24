package goshelf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	"github.com/gorilla/mux"
)

type PathFunction struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
}

const PathPrefix = `/api/` + SchemaVersion + "/"
const BookPath = PathPrefix + `book/`
const CollectionPath = PathPrefix + `collection/`

const applicationJsonContentType = "application/json"

const StatusFailure = "Failure"
const StatusSuccess = "Success"
const CodeFailure int = 400
const CodeSuccess int = 200

func getPathFunctions(cfg *GoshelfConfig) []PathFunction {
	return []PathFunction{
		{
			Path: BookPath,
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					ApiBookFilter(cfg, w, r)
				case http.MethodPost:
					// ApiBookCreate(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
		{
			Path: BookPath + "{id:[0-9]+}",
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					ApiBookGet(cfg, w, r)
				case http.MethodDelete:
					ApiBookRemove(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
		{
			Path: CollectionPath,
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					ApiCollectionCreate(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
		{
			Path: CollectionPath + "{title:[.]+}",
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					ApiCollectionGet(cfg, w, r)
				case http.MethodDelete:
					ApiCollectionDelete(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
	}
}

func checkContentType(contentTypeExpected string, r *http.Request) error {
	val, ok := r.Header["content-type"]
	if !ok {
		return errors.New("content-type header missing")
	}

	lowerCT := strings.ToLower(contentTypeExpected)
	lowerVal := strings.ToLower(val[0])

	if strings.Index(lowerVal, lowerCT) < 0 {
		return errors.New("expected content-type to contain " + contentTypeExpected)
	}

	return nil
}

func ApiCollectionCreate(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	err := checkContentType(applicationJsonContentType, r)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	col := v1.Collection{}

	// Parse the JSON body into object
	err = json.Unmarshal(body, &col)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	// Get all the book ids from within the collection
	bookIds := []int{}
	for _, book := range col.Books {
		bookIds = append(bookIds, book.BookId)
	}

	_, err = cfg.Goshelf.CollectionCreate(&col.Title, col.Books)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	returnGoshelfSuccessWithObject()
}

func ApiCollectionGet(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {

	title := mux.Vars(r)["title"]
	col, err := cfg.Goshelf.CollectionGet(&title)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	ret := map[string]interface{}{
		"collection": col,
	}

	returnGoshelfSuccessWithObject(&ret, w, r)
}

func ApiCollectionDelete(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Should be guaranteed by the regex, but just in case
	idInt64, err := strconv.ParseInt(vars["id"], 10, 32)
	PanicErrorHandler(err)
	id := int(idInt64)

	err = cfg.Goshelf.CollectionRemove(id)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	bookRet := map[string]interface{}{}

	returnGoshelfSuccessWithObject(&bookRet, w, r)
}

func ApiBookRemove(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Should be guaranteed by the regex, but just in case
	idInt64, err := strconv.ParseInt(vars["id"], 10, 32)
	PanicErrorHandler(err)
	id := int(idInt64)

	err = cfg.Goshelf.BookRemove(id)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	bookRet := map[string]interface{}{}

	returnGoshelfSuccessWithObject(&bookRet, w, r)
}

func ApiBookGet(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Should be guaranteed by the regex, but just in case
	idInt64, err := strconv.ParseInt(vars["id"], 10, 32)
	PanicErrorHandler(err)
	id := int(idInt64)

	book, err := cfg.Goshelf.BookGet(id)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	bookRet := map[string]interface{}{
		"book": book,
	}

	returnGoshelfSuccessWithObject(&bookRet, w, r)
}

func ApiBookFilter(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	titleQ := queries.Get("title")
	genreQ := queries.Get("genre")
	editionQ := queries.Get("edition")

	var title *string
	var genre *string
	var edition *int

	if titleQ != "" {
		title = &titleQ
	}

	if genreQ != "" {
		genre = &genreQ
	}

	if editionQ != "" {
		edInt64, err := strconv.ParseInt(editionQ, 10, 32)

		if err != nil {
			errMsg := "edition is not an integer"
			returnGoshelfErrorWithMessage(&errMsg, w, r)
			return
		}

		edInt := int(edInt64)
		edition = &edInt
	}

	books, err := cfg.Goshelf.BookFilter(title, genre, edition)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	bookRet := map[string]interface{}{
		"books": books,
	}

	returnGoshelfSuccessWithObject(&bookRet, w, r)
}

func createResponseObject(status string, code int, metadata *map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":        "sync",
		"status":      status, // "Success", "Error"
		"status_code": code,   // e.g., 400
		"metadata":    metadata,
	}
}

func returnGoshelfErrorWithMessage(msg *string, w http.ResponseWriter, r *http.Request) {
	metadata := map[string]interface{}{
		"message": *msg,
	}

	returnGoshelfResponse(StatusFailure, CodeFailure, &metadata, w, r)
}

func returnGoshelfSuccessWithObject(val *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	returnGoshelfResponse(StatusSuccess, CodeSuccess, val, w, r)
}

func returnGoshelfResponse(status string, code int, metadata *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(createResponseObject(status, code, metadata))
}

func StartServer(cfg GoshelfConfig) {
	r := mux.NewRouter()

	for _, v := range getPathFunctions(&cfg) {
		log.Println("Handling " + v.Path)
		r.HandleFunc(v.Path, v.Function)
	}

	address := cfg.Host + ":" + fmt.Sprint(cfg.Port)

	log.Println("Starting server on " + address)

	err := http.ListenAndServe(address, r)

	if err != nil {
		log.Fatal(err.Error())
	}
}

// func GetCliFuncMap() map[string]func(*GoshelfConfig) {
// 	return cliFuncMap
// }

// func PanicErrorHandler(err error) {
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

// // TODO: automate these based on reflection
// // Creates a book from the cli. Requires a fully configured Goshelfconfig.
// func CliBookCreate(cfg *GoshelfConfig) {
// 	// TODO: implement cli tooling interface
// 	prompt := "\tEnter title: "
// 	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter author's first name: "
// 	aFirst, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter author's first name: "
// 	aLast, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter description (optional): "
// 	desc, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter edition (optional): "
// 	ed, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter genre (optional): "
// 	genre, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter publish date YYYY-MM-dd (optional): "
// 	date, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	book := v1.Book{
// 		Title: *title,
// 		Author: v1.Author{
// 			FirstName: *aFirst,
// 			LastName:  *aLast,
// 		},
// 	}

// 	if *desc != "" {
// 		book.Description = desc
// 	}

// 	if *genre != "" {
// 		book.Genre = genre
// 	}

// 	if *ed != "" {
// 		edInt64, err := strconv.ParseInt(*ed, 10, 32)
// 		edInt := int(edInt64)

// 		if err != nil {
// 			log.Panic("invalid edition (must be integer)")
// 		}

// 		if *ed != "" {
// 			book.Edition = &edInt
// 		}
// 	}

// 	if *date != "" {
// 		pDate, err := time.Parse("2006-01-02", *date)

// 		if err != nil {
// 			log.Panic("invalid time format (must match YYYY-MM-dd)")
// 		}

// 		book.PublishDate = &pDate
// 	}

// 	id, err := cfg.Goshelf.BookCreate(&book)

// 	PanicErrorHandler(err)

// 	fmt.Printf("%v", *id)
// }

// func CliBookGet(cfg *GoshelfConfig) {
// 	prompt := "\tEnter book id: "
// 	idStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	idInt64, err := strconv.ParseInt(*idStr, 10, 32)
// 	PanicErrorHandler(err)

// 	id := int(idInt64)

// 	book, err := cfg.Goshelf.BookGet(id)
// 	PanicErrorHandler(err)

// 	json, err := json.Marshal(book)
// 	PanicErrorHandler(err)

// 	fmt.Print(string(json))
// }

// func CliBookRemove(cfg *GoshelfConfig) {
// 	prompt := "\tEnter book id: "
// 	idStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	idInt64, err := strconv.ParseInt(*idStr, 10, 32)
// 	PanicErrorHandler(err)

// 	id := int(idInt64)

// 	err = cfg.Goshelf.BookRemove(id)
// 	PanicErrorHandler(err)
// }

// func ApiBookFilter(cfg *GoshelfConfig) {
// 	prompt := "\tEnter partial title to search (optional): "
// 	titleStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter edition (optional): "
// 	edStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	prompt = "\tEnter partial genre to search (optional): "
// 	genreStr, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	var title *string
// 	var genre *string
// 	var edition *int

// 	title = nil
// 	genre = nil
// 	edition = nil

// 	if *titleStr != "" {
// 		title = titleStr
// 	}

// 	if *genreStr != "" {
// 		genre = genreStr
// 	}

// 	if *edStr != "" {
// 		edInt64, err := strconv.ParseInt(*edStr, 10, 32)
// 		PanicErrorHandler(err)

// 		editionInt := int(edInt64)
// 		edition = &editionInt
// 	}

// 	books, err := cfg.Goshelf.BookFilter(title, genre, edition)
// 	PanicErrorHandler(err)

// 	for _, book := range books {
// 		json, err := json.Marshal(book)
// 		PanicErrorHandler(err)

// 		fmt.Println(string(json))
// 	}
// }

// func CliCollectionCreate(cfg *GoshelfConfig) {
// 	prompt := "\tEnter collection title: "
// 	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	bookIds := []int{}

// 	for {
// 		prompt := "\tEnter book id (press enter when finished): "
// 		bookId, err := cli.GetIntFromCli(&prompt, os.Stdin, os.Stdout)
// 		PanicErrorHandler(err)

// 		if bookId == nil {
// 			break
// 		}

// 		bookIds = append(bookIds, *bookId)
// 	}

// 	_, err = cfg.Goshelf.CollectionCreate(title, bookIds)
// 	PanicErrorHandler(err)
// }

// func CliCollectionGet(cfg *GoshelfConfig) {
// 	prompt := "\tEnter collection title: "
// 	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	col, err := cfg.Goshelf.CollectionGet(title)
// 	PanicErrorHandler(err)

// 	json, err := json.Marshal(col)
// 	PanicErrorHandler(err)

// 	fmt.Println(string(json))
// }

// func CliCollectionRemove(cfg *GoshelfConfig) {
// 	prompt := "\tEnter collection title: "
// 	title, err := cli.GetCliPrompt(&prompt, os.Stdin)
// 	PanicErrorHandler(err)

// 	err = cfg.Goshelf.CollectionRemove(title)
// 	PanicErrorHandler(err)
// }
