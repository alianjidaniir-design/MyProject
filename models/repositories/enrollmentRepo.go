package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"context"
)

type RegistrationRepository interface {
	// Create method
	Create(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.RegisterStudentRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error)
}
