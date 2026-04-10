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
}

var enrollmentRepos enrollmentReposity = enrollment.GetRepo()
