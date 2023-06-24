package db

import (
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
	RunSpecs(t, "CLI Suite")
}

var _ = Describe("Postgres", func() {

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
		goodBook := v1.Book{
			Author: v1.Author{
				FirstName: "testfirst",
				LastName:  "testlast",
			},
			Title:       "testtitle",
			PublishDate: &now,
			Edition:     &edition,
			Description: &desc,
			Genre:       &genre,
		}

		It("Should save a book", func() {
			pgDb.BookCreate(goodBook)
		})

		// TODO: add more tests
	})

})