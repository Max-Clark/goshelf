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

var _ = Describe("Book Test", func() {

	Context("Connect", func() {
		It("Should successfully connect", func() {
			err := pgDb.Connect()
			Expect(err).To(BeNil())
		})
	})

	Context("Books", func() {
		err := pgDb.Connect()
		Expect(err).To(BeNil())

		Context("BookCreate", func() {
			var bookId *int
			It("Should save a book", func() {
				newBook := BookFactory()

				bookId, err := pgDb.BookCreate(newBook)
				Expect(err).To(BeNil())
				Expect(bookId).ToNot(BeNil())
			})

			AfterEach(func() {
				if bookId != nil {
					pgDb.BookRemove(*bookId)
				}
			})
		})

		Context("BookRead", func() {
			var bookId *int
			var err error

			BeforeEach(func() {
				newBook := BookFactory()
				bookId, err = pgDb.BookCreate(newBook)
				Expect(err).To(BeNil())
				Expect(bookId).ToNot(BeNil())
			})

			It("Should read a book", func() {
				book, err := pgDb.BookGet(*bookId)
				Expect(err).To(BeNil())
				Expect(book).ToNot(BeNil())
			})

			AfterEach(func() {
				if bookId != nil {
					pgDb.BookRemove(*bookId)
				}
			})
		})

		Context("BookRemove", func() {
			var bookId *int
			var err error
			BeforeEach(func() {
				newBook := BookFactory()
				bookId, err = pgDb.BookCreate(newBook)
				Expect(err).To(BeNil())
				Expect(bookId).ToNot(BeNil())
			})

			It("Should delete a book", func() {
				err = pgDb.BookRemove(*bookId)
				Expect(err).To(BeNil())

				book, err := pgDb.BookGet(*bookId)
				Expect(err).To(BeNil())
				Expect(book).To(BeNil())

				// already deleted
				bookId = nil
			})

			It("Should throw error on failed delete", func() {
				err = pgDb.BookRemove(2147483647)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("book not found"))
			})

			AfterEach(func() {
				if bookId != nil {
					pgDb.BookRemove(*bookId)
				}
			})
		})

		Context("BookFilter", func() {
			var booksToSave []v1.Book
			var bookIds []*int
			const bookCardinality = 5

			BeforeEach(func() {
				for i := 0; i < bookCardinality; i++ {
					newBook := BookFactory()
					newBook.Title += fmt.Sprint(i) + "_" + fmt.Sprint(time.Now().UnixMicro())
					*newBook.Genre += fmt.Sprint(i) + "_" + fmt.Sprint(time.Now().UnixMicro())
					edition := i + 1
					newBook.Edition = &edition
					booksToSave = append(booksToSave, *newBook)

					bookId, err := pgDb.BookCreate(newBook)
					Expect(err).To(BeNil())
					Expect(bookId).ToNot(BeNil())

					bookIds = append(bookIds, bookId)
				}
			})

			It("Should filter a book", func() {
				book, err := pgDb.BookFilter(&booksToSave[1].Title, nil, nil)
				Expect(err).To(BeNil())
				Expect(book).ToNot(BeNil())
				Expect(len(book)).To(Equal(1))
			})

			AfterEach(func() {
				for _, bookId := range bookIds {
					if bookId != nil {
						pgDb.BookRemove(*bookId)
					}
				}
			})
		})
		// TODO: add more tests
	})

})

var _ = Describe("Collection Test", func() {

	Context("Connect", func() {
		It("Should successfully connect", func() {
			err := pgDb.Connect()
			Expect(err).To(BeNil())
		})
	})

	Context("Collection", func() {
		err := pgDb.Connect()
		Expect(err).To(BeNil())

		Context("CollectionCreate", func() {
			var booksToSave []v1.Book
			var bookIds []int
			var colTitle string
			const bookCardinality = 5

			BeforeEach(func() {
				bookIds = []int{}
				booksToSave = []v1.Book{}
				for i := 0; i < bookCardinality; i++ {
					newBook := BookFactory()
					newBook.Title += fmt.Sprint(i) + "_" + fmt.Sprint(time.Now().UnixMicro())
					*newBook.Genre += fmt.Sprint(i) + "_" + fmt.Sprint(time.Now().UnixMicro())
					edition := i + 1
					newBook.Edition = &edition
					booksToSave = append(booksToSave, *newBook)

					bookId, err := pgDb.BookCreate(newBook)
					Expect(err).To(BeNil())
					Expect(bookId).ToNot(BeNil())

					bookIds = append(bookIds, *bookId)
				}
			})

			It("Should create a collection", func() {
				colTitle = "collTestCreate" + fmt.Sprint(time.Now().UnixMicro())
				title, err := pgDb.CollectionCreate(&colTitle, bookIds)
				Expect(err).To(BeNil())
				Expect(*title).To(Equal(colTitle))

				collection, err := pgDb.CollectionGet(&colTitle)
				Expect(err).To(BeNil())
				Expect(collection).ToNot(BeNil())
			})

			AfterEach(func() {
				for _, bookId := range bookIds {
					pgDb.BookRemove(bookId)
				}

				pgDb.CollectionRemove(colTitle)
			})
		})

		// TODO: add more tests
	})

})
