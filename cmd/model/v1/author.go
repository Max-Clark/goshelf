package v1

import "time"

type Author struct {
	AuthorId  int `validator:"required,min=1"`
	CreatedTs time.Time
	FirstName string `validator:"required,minLength=1"`
	LastName  string `validator:"required,minLength=1"`
}
