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
