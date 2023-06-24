package v1

type Collection struct {
	Collection_id uint32 `validator:"required,min=1"`
	Title         string `validator:"required,minLength=1"`
	Books         []Book `validator:"required"`
}
