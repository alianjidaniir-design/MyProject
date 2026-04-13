package teacherSchema

type InformationSchema struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
}
