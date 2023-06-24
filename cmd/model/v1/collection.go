package v1

import "time"

type Collection struct {
	Title     string    `validator:"required,minLength=1" json:"title"`
	CreatedTs time.Time `json:"createdTs,omitempty"`
	Books     []Book    `validator:"required" json:"books"`
}
