package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations"
	"context"
)

type RegisterRepository interface {
	CreateRegistration(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.RegisterStudentRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.DeleteStudentResponse, errStr string, code int, err error)
}

var RegistrationRepo RegisterRepository = Registrations.GetRepo()
