package v1

type Collection struct {
	Title string `validator:"required,minLength=1"`
	Books []Book `validator:"required"`
}
