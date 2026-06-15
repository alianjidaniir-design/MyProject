package dataModels

import (
	"time"
)

type Offering struct {
	Row            int64      `json:"row"`
	GroupNumber    int        `json:"group_number"`
	CourseNumber   int64      `json:"course_number"`
	TeacherID      int64      `json:"teacher_id"`
	Capacity       int        `json:"capacity"`
	EnrolledCount  int64      `json:"enrolled_count"`
	IsActive       bool       `json:"is_active"`
	Reservation    int64      `json:"reservation"`
	TermID         int64      `json:"term_id"`
	Week           string     `json:"week"`
	Day            string     `json:"day"`
	ClassStartTime string     `json:"class_start_time"`
	ClassEndTime   string     `json:"class_end_time"`
	ExamStartTime  *time.Time `json:"exam_start_time"`
	ExamFinishTime *time.Time `json:"exam_finish_time"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ListOfferings struct {
	Row              int64      `json:"row"`
	GroupNumber      int        `json:"group_number"`
	CourseNumber     int64      `json:"course_number"`
	Title            string     `json:"title"`
	Unit             int        `json:"unit"`
	TeacherName      string     `json:"teacher_name"`
	TeacherLastName  string     `json:"teacher_last_name"`
	College          string     `json:"college"`
	EducationalGroup string     `json:"educational_group"`
	Capacity         int        `json:"capacity"`
	EnrolledCount    int64      `json:"enrolled_count"`
	IsActive         bool       `json:"is_active"`
	Reservation      int64      `json:"reservation"`
	Week             string     `json:"week"`
	Day              string     `json:"day"`
	ClassStartTime   string     `json:"class_start_time"`
	ClassEndTime     string     `json:"class_end_time"`
	ExamStartTime    *time.Time `json:"exam_start_time"`
	ExamEndTime      *time.Time `json:"exam_end_time"`
	Term             int        `json:"term"`
	Year             int        `json:"year"`
}

type DetailClasses struct {
	OfferingGroupNumber int    `json:"OfferingGroupNumber"`
	CourseNumber        int    `json:"CourseNumber"`
	Title               string `json:"title"`
	Capacity            int    `json:"capacity"`
	EnrolledCount       int    `json:"EnrolledCount"`
	Week                string `json:"week"`
	WorkDay             string `json:"workDay"`
	ClassStartTime      string `json:"ClassStartTime"`
	ClassEndTime        string `json:"ClassEndTime"`
	ExamStartTime       string `json:"ExamStartTime"`
	ExamFinishTime      string `json:"ExamFinishTime"`
	DepartmentID        int64  `json:"DepartmentID"`
}

type TermClasses struct {
	Term    int `json:"term"`
	Year    int `json:"year"`
	Classes []DetailClasses
}
