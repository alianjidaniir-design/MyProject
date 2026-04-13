package dataModels

import "time"

type Teacher struct {
	ID             int64     `gorm:"column:ID" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	LastName       string    `gorm:"column:last_name" json:"last_name"`
	Email          string    `gorm:"column:email" json:"email"`
	Phone          string    `gorm:"column:phone" json:"phone"`
	WorkExperience string    `gorm:"column:work_experience" json:"work_experience"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
