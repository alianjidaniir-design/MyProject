package dataSources

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
	"context"
)

type RegistrationDS interface {
	RegistrationsStudent(ctx context.Context, req registrationSchema.RegisterStudentRequest) (res dataModels.Registration, err error)
	GetRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error)
	UpdateRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error)
	DeleteRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error)
	ListAllRegisterStudent(ctx context.Context, req registrationSchema.SelectPageRegisteredStudentsRequest) (res []dataModels.Registration, total int, err error)
	CancelRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error)
}
