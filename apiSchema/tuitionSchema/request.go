package tuitionSchema

type CreateTuition struct {
	StudentID     int64 `json:"student_id"`
	OfferingRow   int64 `json:"offering_row"`
	FixedTuition  int   `json:"fixed_tuition"`
	CourseTuition int   `json:"course_tuition"`
	ExtraOption   int   `json:"extra_option"`
}

type UpdateTuition struct {
	Row           int64 `json:"row"`
	CourseTuition int   `json:"course_tuition"`
	ExtraOption   int   `json:"extra_option"`
}
