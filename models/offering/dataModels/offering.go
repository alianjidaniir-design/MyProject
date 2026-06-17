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
	Row              int64      `json:"Row"`
	GroupNumber      int        `json:"GroupNumber"`
	CourseNumber     int64      `json:"CourseNumber"`
	Title            string     `json:"Title"`
	Unit             int        `json:"Unit"`
	TeacherName      string     `json:"TeacherName"`
	TeacherLastName  string     `json:"TeacherLastName"`
	College          string     `json:"College"`
	EducationalGroup string     `json:"EducationalGroup"`
	Capacity         int        `json:"Capacity"`
	EnrolledCount    int64      `json:"EnrolledCount"`
	IsActive         bool       `json:"IsActive"`
	Reservation      int64      `json:"Reservation"`
	Week             string     `json:"Week"`
	Day              string     `json:"Day"`
	ClassStartTime   string     `json:"ClassStartTime"`
	ClassEndTime     string     `json:"ClassEndTime"`
	ExamStartTime    *time.Time `json:"ExamStartTime"`
	ExamEndTime      *time.Time `json:"ExamEndTime"`
	Term             int        `json:"Term"`
	Year             int        `json:"Year"`
}

type DetailClasses struct {
	OfferingGroupNumber int     `json:"OfferingGroupNumber"`
	CourseNumber        int     `json:"CourseNumber"`
	Title               string  `json:"Title"`
	Capacity            int     `json:"Capacity"`
	EnrolledCount       int     `json:"EnrolledCount"`
	Week                string  `json:"Week"`
	WorkDay             string  `json:"WorkDay"`
	ClassStartTime      string  `json:"ClassStartTime"`
	ClassEndTime        string  `json:"ClassEndTime"`
	ExamStartTime       *string `json:"ExamStartTime"`
	ExamFinishTime      *string `json:"ExamFinishTime"`
	DepartmentID        int64   `json:"DepartmentID"`
}

type TermClasses struct {
	Term    int `json:"Term"`
	Year    int `json:"Year"`
	Classes []DetailClasses
}
