package v1

import "time"

type Collection struct {
	Title     string `validator:"required,minLength=1"`
	CreatedTs time.Time
	Books     []Book `validator:"required"`
}
