package dataModels

import "time"

type Teacher struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	LastName       string     `json:"last_name"`
	RoleName       string     `json:"role_name"`
	NationalCode   string     `json:"national_code"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	WorkExperience string     `json:"work_experience"`
	Password       string     `json:"password"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
