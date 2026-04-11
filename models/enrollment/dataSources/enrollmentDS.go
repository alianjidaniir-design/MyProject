package dataSources

import (
	"MyProject/apiSchema/enrollmentSchema"
	EnrollmentDataModel "MyProject/models/enrollment/dataModels"
	"context"
)

type EnrollmentDS interface {
	EnrollStudent(ctx context.Context, req enrollmentSchema.EnrollmentRequest) (res EnrollmentDataModel.Enrollment, err error)
	CancelEnrollment(ctx context.Context, req enrollmentSchema.CancelEnrollmentRequest) (res EnrollmentDataModel.Enrollment, err error, result string)
	ListEnrollment(ctx context.Context, req enrollmentSchema.ListEnrollmentsRequest) (res []EnrollmentDataModel.Enrollment, code int64, err error)
	ListStudentCourses(ctx context.Context, req enrollmentSchema.ListStudentCoursesRequest) (res EnrollmentDataModel.Enrollment, err error)
}
