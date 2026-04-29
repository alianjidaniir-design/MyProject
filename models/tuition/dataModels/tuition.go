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
