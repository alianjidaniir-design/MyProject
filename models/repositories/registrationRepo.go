package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations"
	"context"
)

type RegisterRepository interface {
	CreateRegistration(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.RegisterStudentRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error)
}

var RegistrationRepo RegisterRepository = Registrations.GetRepo()
