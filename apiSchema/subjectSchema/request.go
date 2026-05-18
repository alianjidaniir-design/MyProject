package subjectSchema

type CreateSubject struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty,max=127"`
}

type GetSubject struct {
	ID int64 `json:"id" validate:"required"`
}


type Pagination struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"required"`
}