package dataModels

import "time"

type Course struct {
	ID           int64     `json:"id"`
	CourseNumber string    `json:"course_number"`
	Title        string    `json:"title"`
	Unit         int       `json:"unit"`
	DepartmentID string    `json:"department_ID"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
