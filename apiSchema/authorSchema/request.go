package authorSchema

type CreateAuthor struct {
	FirstName string `json:"first_name" val:"required,max=31"`
	LastName  string `json:"last_name" val:"required,max=63"`
	BirthYear int    `json:"birth_year" val:"required,min=1300,max=9999"`
}

type GetAuthor struct {
	ID int64 `json:"id"`
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
