package studentSchema

type SignUpStudent struct {
	Name         string  `json:"name" validate:"required,max=31"`
	Family       string  `json:"family" validate:"required,max=63"`
	Phone        string  `json:"phone" validate:"required,numeric,len=11"`
	NationalCode string  `json:"national_code" validate:"required,numeric,len=10"`
	Major        string  `json:"major" validate:"required"`
	StudentCode  string  `json:"student_code" validate:"required,numeric,len=9"`
	TermID       int64   `json:"term_id" validate:"required"`
	Level        string  `json:"level" validate:"required"`
	UserName     *string `json:"user_name" validate:"omitempty,len=9,numeric"`
	Password     string  `json:"password" validate:"required,len=10"`
	RoleName     string  `json:"role_name"`
}
type LoginStudent struct {
	UserName string `json:"user_name" validate:"required,len=9,numeric"`
	Password string `json:"password" validate:"required,len=10"`
}

type ListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetRequest struct {
	ID int64 `json:"ID"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"ID"`
	UserName string `json:"user_name" validate:"required,len=9,numeric"`
}
type DeleteRequest struct {
	ID int64 `json:"ID"`
}

type SoftDeleteRequest struct {
	ID int64 `json:"ID"`
}

type LogoutRequest struct {
	IsLogout bool `json:"is_logout"`
}
