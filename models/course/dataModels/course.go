package dataModels

import "time"

type Course struct {
	ID         int64     `gorm:"column:ID" json:"id"`
	CourseCode string    `gorm:"column:course_code" json:"course_code"`
	Title      string    `gorm:"column:title" json:"title"`
	Capacity   int       `gorm:"column:capacity" json:"capacity"`
	EnrolledAt int       `gorm:"column:enrolled_at" json:"enrolled_at"`
	IsActive   bool      `gorm:"column:isActive" json:"isActive"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
