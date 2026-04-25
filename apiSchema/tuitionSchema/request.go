package tuitionSchema

type CreateTuition struct {
	StudentID     int64 `json:"student_id"`
	CourseID      int64 `json:"course_id"`
	FixedTuition  int   `json:"fixed_tuition"`
	CourseTuition int   `json:"course_tuition"`
	ExtraOption   int   `json:"extra_option"`
	DebitAmount   int   `json:"debit_amount"`
	CreditAmount  int   `json:"credit_amount"`
}
