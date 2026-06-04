package dataModels

import "time"

type Registration struct {
	ID           int64      `json:"id"`
	StudentID    int64      `json:"student_id"`
	CourseNumber int64      `json:"course_number"`
	OfferingRow  int64      `json:"offering_row"`
	Status       string     `json:"status"`
	Registrar    string     `json:"registrar"`
	EnrolledAt   time.Time  `json:"enrolled_at"`
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

type DetailClassScheduleR struct {
	OfferingRow     int64     `json:"offering_row"`
	OfferingNumber  int       `json:"offering_number"`
	CourseNumber    int       `json:"course_number"`
	Title           string    `json:"title"`
	Unit            int       `json:"unit"`
	TeacherID       int64     `json:"teacher_id"`
	TeacherName     string    `json:"teacher_name"`
	TeacherLastName string    `json:"teacher_last_name"`
	ClassStartTime  time.Time `json:"class_start_time"`
	ClassEndTime    time.Time `json:"class_end_time"`
}

type TermClassSchedules struct {
	Term       int `json:"term"`
	Year       int `json:"year"`
	classes    []DetailClassScheduleR
	TotalUnits int `json:"total_units"`
}
