package v1

import "time"

type Book struct {
	BookId       uint32     `validator:"required,min=1"`
	Author       Author     `validator:"required"`
	Title        string     `validator:"required,minLength=1"`
	Publish_date *time.Time `validator:"optional"`
	Edition      *int       `validator:"optional,min=1"`
	Description  *string    `validator:"optional,minLength=1"`
	Genre        *string    `validator:"optional"`
}
