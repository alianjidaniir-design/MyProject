package teacherSchema

type InformationSchema struct {
	Name           string `json:"name" validate:"required,max=20,alpha"`
	LastName       string `json:"last_name" validate:"required,max=63,alpha"`
	NationalCode   string `json:"national_code" validate:"required,len=10,numeric"`
	Email          string `json:"email" validate:"required,email"`
	Phone          string `json:"phone" validate:"required,len=11,numeric"`
	WorkExperience string `json:"work_experience" validate:"omitempty,max=255"`
	Password       string `json:"password" validate:"required,len=11,numeric"`
}

type LoginTeacherRequest struct {
	NationalCode string `json:"national_code" validate:"required,len=10,numeric"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,len=11,numeric"`
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
