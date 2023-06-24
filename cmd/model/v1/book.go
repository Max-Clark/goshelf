package v1

import "time"

type Book struct {
	BookId      int        `validator:"required,min=1"`
	Author      Author     `validator:"required"`
	Title       string     `validator:"required,minLength=1"`
	PublishDate *time.Time `validator:"optional"`
	Edition     *int       `validator:"optional,min=1"`
	Description *string    `validator:"optional,minLength=1"`
	Genre       *string    `validator:"optional"`
}

func GetEmptyBook() *Book {
	res := Book{
		BookId:      0,
		Title:       "",
		PublishDate: &time.Time{},
		Edition:     nil,
		Description: nil,
		Genre:       nil,
		Author: Author{
			AuthorId:  0,
			FirstName: "",
			LastName:  "",
		},
	}

	return &res
}
