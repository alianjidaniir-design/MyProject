package dataSources

import (
	"MyProject/apiSchema/enrollmentSchema"
	EnrollmentDataModel "MyProject/models/enrollment/dataModels"
	"context"
)

type EnrollmentDS interface {
	EnrollStudent(ctx context.Context, req enrollmentSchema.EnrollmentRequest) (res EnrollmentDataModel.Enrollment, err error)
}
