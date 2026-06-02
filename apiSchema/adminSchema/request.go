package adminSchema

type InformationSchema struct {
	Username string  `json:"user_name" validate:"required,max=15"`
	Password string  `json:"password" validate:"required,len=11"`
	Name     string  `json:"name" validate:"required,max=31"`
	Family   string  `json:"family" validate:"required,max=63"`
	Email    *string `json:"email" validate:"omitempty,email,max=127"`
	RoleName string  `json:"role_name" validate:"required"`
}

type LoginAdminRequest struct {
	Username string `json:"user_name" validate:"required,max=15"`
	Password string `json:"password" validate:"required,len=11"`
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

type LogoutSchema struct {
	IsLogout bool `json:"is_logout" validate:"required"`
}
