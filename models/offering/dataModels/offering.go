package dataModels

import "time"

type Offering struct {
	Row            int64     `json:"row"`
	GroupNumber    int       `json:"group_number"`
	CourseID       int64     `json:"course_id"`
	TeacherID      int64     `json:"teacher_id"`
	Capacity       int       `json:"capacity"`
	EnrolledCount  int64     `json:"enrolled_count"`
	IsActive       bool      `json:"is_active"`
	Reservation    int64     `json:"reservation"`
	TermID         int64     `json:"term_id"`
	ClassStartTime time.Time `json:"class_start_time"`
	ClassEndTime   time.Time `json:"class_end_time"`
	ExamStartTime  time.Time `json:"exam_start_time"`
	ExamEndTime    time.Time `json:"exam_end_time"`
}
