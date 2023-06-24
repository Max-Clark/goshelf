package v1

import "time"

type Book struct {
	BookId      int        `validator:"required,min=1" json:"bookId"`
	Author      Author     `validator:"required" json:"author"`
	CreatedTs   time.Time  `json:"createdTs"`
	Title       string     `validator:"required,minLength=1" json:"title"`
	PublishDate *time.Time `validator:"optional" json:"publishDate,omitempty"`
	Edition     *int       `validator:"optional,min=1" json:"edition,omitempty"`
	Description *string    `validator:"optional,minLength=1" json:"description,omitempty"`
	Genre       *string    `validator:"optional" json:"genre,omitempty"`
}
