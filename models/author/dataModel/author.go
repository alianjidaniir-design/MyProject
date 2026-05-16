package dataModel

type Author struct {
	ID        int64  `json:"author_id" validate:"required"`
	FirstName string `json:"first_name" validate:"required , len=31"`
	LastName  string `json:"last_name" validate:"required,len=63"`
	BirthYear int    `json:"birth_year" validate:"required,min=1300,max=9999"`
}
