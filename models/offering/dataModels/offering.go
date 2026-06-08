package dataModels

import "time"

type Offering struct {
	Row            int64     `json:"row"`
	GroupNumber    int       `json:"group_number"`
	CourseNumber   int64     `json:"course_number"`
	TeacherID      int64     `json:"teacher_id"`
	Capacity       int       `json:"capacity"`
	EnrolledCount  int64     `json:"enrolled_count"`
	IsActive       bool      `json:"is_active"`
	Reservation    int64     `json:"reservation"`
	TermID         int64     `json:"term_id"`
	Week           string    `json:"week"`
	Day            string    `json:"day"`
	ClassStartTime string    `json:"class_start_time"`
	ClassEndTime   string    `json:"class_end_time"`
	ExamStartTime  time.Time `json:"exam_start_time"`
	ExamEndTime    time.Time `json:"exam_end_time"`
}
