package dataModels

import "time"

type Enrollment struct {
	ID         int64      `json:"id"`                    // Primary key, auto-incrementing
	StudentID  int64      `json:"student_id"`            // Foreign key to the student
	CourseID   int64      `json:"course_id"`             // Foreign key to the course
	Status     string     `json:"status"`                // Enrollment status ('enrolledd', 'canceled')
	EnrolledAt time.Time  `json:"enrolled_at"`           // Timestamp of enrollment
	CanceledAt *time.Time `json:"canceled_at,omitempty"` // Timestamp of cancellation (can be null)
	CreatedAt  time.Time  `json:"created_at"`            // Timestamp when the record was created
	UpdatedAt  time.Time  `json:"updated_at"`            // Timestamp when the record was last updated
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`  // Timestamp for soft delete (can be null)
}
