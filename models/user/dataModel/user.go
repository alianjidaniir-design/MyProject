package dataModel

import "time"

type User struct {
	ID        int64     `gorm:"column:ID" json:"ID"`
	Code      string    `gorm:"column:code" json:"code"`
	Name      string    `gorm:"column:name" json:"name"`
	Family    string    `gorm:"column:family" json:"family"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Major     string    `gorm:"column:major" json:"major"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
