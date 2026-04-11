package enrollmentSchema

type EnrollmentRequest struct {
	StudentID int64 `json:"student_id"`
	CourseID  int64 `json:"course_id"`
}

type CancelEnrollmentRequest struct {
	ID int64 `json:"id"`
}

type ListEnrollmentsRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ListStudentCoursesRequest struct {
	StudentID int64  `json:"student_id"`
	Status    string `json:"status"`
}
