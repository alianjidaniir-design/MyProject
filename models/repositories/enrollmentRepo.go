package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/enrollmentSchema"
	"MyProject/models/enrollment"
	"context"
)

type enrollmentReposity interface {
	// Create method
	Create(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.EnrollmentRequest]) (res enrollmentSchema.EnrollmentResponse, errStr string, code int, err error)
	Cancel(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.CancelEnrollmentRequest]) (res enrollmentSchema.DeactivateEnrollmentResponse, errStr string, code int, err error)
	ListEnrollment(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.ListEnrollmentsRequest]) (res enrollmentSchema.ListEnrollmentResponse, errStr string, code int, err error)
	ListStudentCourse(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.ListStudentCoursesRequest]) (res enrollmentSchema.ListStudentCoursesResponse, errStr string, code int, err error)
}

var EnrollmentRepos enrollmentReposity = enrollment.GetRepo()
