package dataModels

import "time"

type Course struct {
	ID            int64     `json:"id"`
	CourseCode    string    `json:"course_code"`
	Title         string    `json:"title"`
	TeacherID     int64     `json:"teacher_id"`
	Credits       int       `json:"credit"`
	Capacity      int       `json:"capacity"`
	EnrolledCount int       `json:"enrolled_count"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}
