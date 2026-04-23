package dataModels

import "time"

type Registration struct {
	ID          int64     `json:"id"`
	StudentID   int64     `json:"student_id"`
	CourseID    int64     `json:"course_id"`
	OfferingRow int64     `json:"offering_row"`
	Status      string    `json:"status"`
	EnrolledAt  time.Time `json:"enrolled_at"`
	CanceledAt  time.Time `json:"canceled_at "`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type Student struct {
	StudentID int64  `json:"student_id"`
	Status    string `json:"status"`
}

type Offering struct {
	OfferingRow int64  `json:"offering_row"`
	Status      string `json:"status"`
}
