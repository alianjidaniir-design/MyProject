package dataModels

import "time"

type Admins struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Family    string    `json:"family"`
	Email     *string   `json:"email"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}
