package dataModel

type Translator struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthYear int    `json:"birth_year"`
}
