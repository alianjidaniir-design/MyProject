package dataModels

import "time"

type Tuition struct {
	Row           int64     `json:"row"`
	StudentID     int64     `json:"student_id"`
	OfferingID    int64     `json:"offering_id"`
	FixedTuition  int       `json:"fixed_tuition"`
	CourseTuition int       `json:"course_tuition"`
	ExtraOption   int       `json:"extra_option"`
	DebitAmount   int       `json:"debit_amount"`
	CreditAmount  int       `json:"credit_amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

type StudentsDebit struct {
	StudentID    int64 `json:"student_id"`
	TotalTuition int   `json:"total_tuition"`
}

type TuitionStudent struct {
	StudentID     int64  `json:"student_id"`
	StudentCode   string `json:"student_code"`
	StudentName   string `json:"student_name"`
	StudentFamily string `json:"student_family"`
	Major         string `json:"major"`
	CourseID      int64  `json:"course_id"`
	CourseTitle   string `json:"course_title"`
	CourseNumber  int    `json:"course_number"`
	Unit          string `json:"unit"`
	FixedTuition  int    `json:"fixed_tuition"`
	CourseTuition int    `json:"course_tuition"`
	ExtraOption   int    `json:"extra_option"`
	DebitAmount   int    `json:"debit_amount"`
	CreditAmount  int    `json:"credit_amount"`
}
