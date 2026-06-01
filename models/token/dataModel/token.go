package dataModel

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	UserID   int64  `json:"user_id"`
	RoleName string `json:"role_name"`
	Scope    string `json:"scope"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	UserID    int64     `gorm:"column:user_id"`
	RoleName  string    `gorm:"column:role_name"`
	Token     string    `gorm:"column:token,unique;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	RevokedAt bool      `gorm:"default:false"`
}
