package postgresql

import (
	"fmt"
	"testing"
	"time"

	db "github.com/Max-Clark/goshelf/cmd/db"
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var pgDb = PgDb{
	SchemaVersion: "v1",
	Config: db.ConnectionConfig{
		Host:     "0.0.0.0",
		Port:     5432,
		User:     "postgres",         // TODO: Grab from ENV or Arg
		Password: "mysecretpassword", // TODO: Grab from ENV or Arg
		DbName:   "postgres",
		SslMode:  "disable",
	},
}

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Suite")
}

func BookFactory() *v1.Book {
	now := time.Now()
	desc := "bookDesc"
	genre := "bookGenre"
	edition := 1
	newAuthor := v1.Author{
		FirstName: "testfirst" + fmt.Sprint(time.Now().UnixMilli()),
		LastName:  "testlast" + fmt.Sprint(time.Now().UnixMilli()),
	}

	return &v1.Book{
		Author:      newAuthor,
		Title:       "testtitle",
		PublishDate: &now,
		Edition:     &edition,
		Description: &desc,
		Genre:       &genre,
	}
}
