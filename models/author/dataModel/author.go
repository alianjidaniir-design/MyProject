package dataModel

type Author struct {
	ID        int64  `json:"author_id" val:"required"`
	FirstName string `json:"first_name" val:"required , len=31"`
	LastName  string `json:"last_name" val:"required,len=63"`
	BirthYear int    `json:"birth_year" val:"required,min=1300,max=9999"`
}
