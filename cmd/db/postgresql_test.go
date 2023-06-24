package db

import (
	"fmt"
	"testing"
	"time"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var pgDb = PgDb{
	SchemaVersion: "v1",
	Config: ConnectionConfig{
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

var _ = Describe("Postgres Test", func() {

	Context("Connect", func() {
		It("Should successfully connect", func() {
			err := pgDb.Connect()
			Expect(err).To(BeNil())
		})
	})

	Context("Books", func() {
		err := pgDb.Connect()
		Expect(err).To(BeNil())

		now := time.Now()
		desc := "What a cool book"
		genre := "horror"
		edition := 1
		newAuthor := v1.Author{
			FirstName: "testfirst" + fmt.Sprint(time.Now().UnixMilli()),
			LastName:  "testlast" + fmt.Sprint(time.Now().UnixMilli()),
		}

		newAuthorBook := v1.Book{
			Author:      newAuthor,
			Title:       "testtitle",
			PublishDate: &now,
			Edition:     &edition,
			Description: &desc,
			Genre:       &genre,
		}

		It("Should save and return a book with a new author", func() {
			savedId, err := pgDb.BookCreate(&newAuthorBook)
			Expect(err).To(BeNil())
			Expect(savedId).ToNot(BeNil())

			book, err := pgDb.BookGet(*savedId)
			Expect(err).To(BeNil())
			Expect(book).ToNot(BeNil())

			Expect(book.Author.FirstName).To(Equal(newAuthor.FirstName))
			Expect(book.Author.LastName).To(Equal(newAuthor.LastName))
		})

		// TODO: add more tests
	})

})
