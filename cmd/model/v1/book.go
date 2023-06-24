package v1

import "time"

type Book struct {
	BookId      int    `validator:"required,min=1"`
	Author      Author `validator:"required"`
	CreatedTs   time.Time
	Title       string     `validator:"required,minLength=1"`
	PublishDate *time.Time `validator:"optional"`
	Edition     *int       `validator:"optional,min=1"`
	Description *string    `validator:"optional,minLength=1"`
	Genre       *string    `validator:"optional"`
}
