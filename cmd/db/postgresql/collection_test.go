package postgresql

import (
	"fmt"
	"time"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

				pgDb.CollectionRemove(&colTitle)
			})
		})

		// TODO: add more tests
	})

})
