package dataModel

import (
	"MyProject/models/role/dataModel"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Student struct {
	ID           int64           `gorm:"column:ID" json:"ID"`
	Name         string          `gorm:"column:name" json:"name"`
	Family       string          `gorm:"column:family" json:"family"`
	Phone        string          `gorm:"column:phone" json:"phone"`
	NationalCode string          `gorm:"column:national_code" json:"national_code"`
	Major        string          `gorm:"column:major" json:"major"`
	StudentCode  string          `gorm:"column:student_code" json:"student_code"`
	UserName     string          `gorm:"column:user_name" json:"user_name"`
	Password     string          `gorm:"column:password" json:"password"`
	RoleID       int64           `gorm:"column:role_id" json:"role_id"`
	Role         *dataModel.Role `json:"role"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time      `gorm:"column:deleted_at" json:"deleted_at"`
}

type AccessToken struct {
	StudentID int64  `json:"student_id"`
	RoleID    string `json:"role_id"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	StudentID int64     `gorm:"column:student_id"`
	Token     string    `gorm:"column:token,unique;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	RevokedAt bool      `gorm:"default:false"`
}
