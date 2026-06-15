package dataModels

import "time"

type Registration struct {
	ID           int64      `json:"id"`
	StudentID    int64      `json:"student_id"`
	CourseNumber int64      `json:"course_number"`
	OfferingRow  int64      `json:"offering_row"`
	Status       string     `json:"status"`
	Registrar    string     `json:"registrar"`
	RegisteredAt time.Time  `json:"registered_at"`
	CanceledAt   *time.Time `json:"canceled_at "`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ListSelectOfferingResponse struct {
	ID           int64  `json:"id"`
	StudentID    int64  `json:"student_id"`
	CourseNumber int64  `json:"course_number"`
	OfferingRow  int64  `json:"offering_row"`
	Status       string `json:"status"`
	Err          string `json:"error"`
}

type Student struct {
	StudentID int64  `json:"student_id"`
	Status    string `json:"status"`
}

type Offering struct {
	OfferingRow int64  `json:"offering_row"`
	Status      string `json:"status"`
}

type DetailClassScheduler struct {
	OfferingRow         int64  `json:"offering_row"`
	OfferingGroupNumber int    `json:"offering_group_number"`
	CourseNumber        int    `json:"course_number"`
	Title               string `json:"title"`
	Unit                int    `json:"unit"`
	TeacherName         string `json:"teacher_name"`
	TeacherLastName     string `json:"teacher_last_name"`
	ClassStartTime      string `json:"class_start_time"`
	ClassEndTime        string `json:"class_end_time"`
}

type TermClassSchedules struct {
	Term       int `json:"term"`
	Year       int `json:"year"`
	Classes    []DetailClassScheduler
	TotalUnits int `json:"total_units"`
}
