package teacherSchema

type InformationSchema struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
}

type PaginationSchema struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetTeacherSchema struct {
	ID int64 `json:"id"`
}

type SelectTeacherSchema struct {
	ID int64 `json:"id"`
}
