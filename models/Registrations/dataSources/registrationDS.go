package dataSources

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
	"context"
)

type RegistrationDS interface {
	RegistrationsStudent(ctx context.Context, req registrationSchema.RegisterStudentRequest) (res dataModels.Registration, err error)
}
