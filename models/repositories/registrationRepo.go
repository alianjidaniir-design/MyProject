package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations"
	"context"

	"github.com/gofiber/fiber/v2"
)

type RegisterRepository interface {
	CreateRegistration(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.RegisterStudentRequest], c *fiber.Ctx) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.DetailStudentResponse, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.DetailStudentResponse, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest], c *fiber.Ctx) (res registrationSchema.DeleteStudentResponse, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.SelectPageRegisteredStudentsRequest]) (res registrationSchema.ListStudentsResponse, errStr string, code int, err error)
	Cancel(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest], c *fiber.Ctx) (res registrationSchema.CancelRegistrationResponse, errStr string, code int, err error)
	ListStudents(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListStudentsRequest]) (res registrationSchema.ListStudentResponse, errStr string, code int, err error)
	ListOfferings(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListOfferingRequest]) (res registrationSchema.ListOfferingResponse, errStr string, code int, err error)
	ListClassesStudent(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.Pages], c *fiber.Ctx) (res registrationSchema.ClassSchedule, errStr string, code int, err error)
	ListClassesTeacher(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.Pages], c *fiber.Ctx) (res registrationSchema.ClassesTeacher, errStr string, code int, err error)
}

var RegistrationRepo RegisterRepository = Registrations.GetRepo()
