package studentSchema

type SignUpStudent struct {
	Name         string  `json:"name" validate:"required,max=31"`
	Family       string  `json:"family" validate:"required,max=63"`
	Phone        string  `json:"phone" validate:"required,numeric,len=11"`
	NationalCode string  `json:"national_code" validate:"required,numeric,len=10"`
	Major        string  `json:"major" validate:"required"`
	StudentCode  string  `json:"student_code" validate:"required,numeric,len=9"`
	UserName     *string `json:"user_name" validate:"omitempty,len=9,numeric"`
	Password     string  `json:"password" validate:"required,len=10"`
	RoleID       int64   `json:"role_id"`
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
	ID int64 `json:"ID"`
}
type DeleteRequest struct {
	ID int64 `json:"ID"`
}

type SoftDeleteRequest struct {
	ID int64 `json:"ID"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	IsLogout bool `json:"is_logout"`
}
