package authorSchema

type CreateAuthor struct {
	FirstName string `json:"first_name" validate:"required,max=31"`
	LastName  string `json:"last_name" validate:"required,max=63"`
	BirthYear int    `json:"birth_year" validate:"required,min=1300,max=9999"`
}

type GetAuthor struct {
	ID int64 `json:"id" validate:"required"`
}

type Pagination struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"required"`
}
