package goshelf

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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

type CollectionCreateApiStruct struct {
	Title   string `json:"title"`
	BookIds []int  `json:"bookIds"`
}

func ApiCollectionCreate(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	err := checkContentType(applicationJsonContentType, r)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	col := CollectionCreateApiStruct{}

	// Parse the JSON body into object
	err = json.Unmarshal(body, &col)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	_, err = cfg.Goshelf.CollectionCreate(&col.Title, col.BookIds)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	returnGoshelfSuccessWithNoObject(w, r)
}

func ApiCollectionGet(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {

	title := mux.Vars(r)["title"]
	col, err := cfg.Goshelf.CollectionGet(&title)

	if err != nil {
		errMsg := err.Error()
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	if col == nil {
		errMsg := "not found"
		returnGoshelfErrorWithMessage(&errMsg, w, r)
		return
	}

	ret := map[string]interface{}{
		"collection": col,
	}

	returnGoshelfSuccessWithObject(&ret, w, r)
}

func ApiCollectionDelete(cfg *GoshelfConfig, w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	err := cfg.Goshelf.CollectionRemove(&title)

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
