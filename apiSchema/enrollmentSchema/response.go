package enrollmentSchema

import "MyProject/models/enrollment/dataModels"

type EnrollmentResponse struct {
	Enrollment dataModels.Enrollment `json:"enrollment"`
}

type DeactivateEnrollmentResponse struct {
	Enrollment dataModels.Enrollment `json:"enrollment"`
	Result     string                `json:"result"`
}

type ListEnrollmentResponse struct {
	Enrollments []dataModels.Enrollment `json:"enrollments"`
	TotalCount  int64                   `json:"total_count"`
}

type ListStudentCoursesResponse struct {
	Enrollments []dataModels.Enroll `json:"enrollments"`
}

type ListCourseStudentsResponse struct {
	Course []dataModels.Course `json:"enrollments"`
}
