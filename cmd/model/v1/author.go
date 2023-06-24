package v1

import "time"

type Author struct {
	AuthorId  int       `validator:"required,min=1" json:"authorId,omitempty"`
	CreatedTs time.Time `json:"createdTs,omitempty"`
	FirstName string    `validator:"required,minLength=1" json:"firstName,omitempty"`
	LastName  string    `validator:"required,minLength=1" json:"lastName,omitempty"`
}
