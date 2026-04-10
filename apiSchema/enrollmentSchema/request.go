package enrollmentSchema

type EnrollmentRequest struct {
	StudentID int64 `json:"student_id"`
	CourseID  int64 `json:"course_id"`
}
