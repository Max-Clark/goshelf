package v1

import "time"

type Author struct {
	AuthorId  uint32 `validator:"required,min=1"`
	Timestamp time.Time
	FirstName string `validator:"required,minLength=1"`
	LastName  string `validator:"required,minLength=1"`
}
