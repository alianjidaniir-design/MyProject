package route

import (
	. "MyProject/controllers/registration"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var registerRoute = map[string]string{
	"registrationCreate":        "registration/create",
	"registrationGet":           "registration/get",
	"registrationUpdate":        "registration/update",
	"registrationDelete":        "registration/delete",
	"registrationList":          "registration/list",
	"registrationCancel":        "registration/cancel",
	"registrationListStudent":   "registration/student",
	"registrationListOffering":  "registration/offering",
	"registrationListMyClasses": "registration/myclasses",
}

func SetupRegistrationRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(registerRoute["registrationCreate"], authz.RequirePermission(permissions.CreateRegister), Create)
	api.Post(registerRoute["registrationGet"], authz.RequirePermission(permissions.ViewRegister), Get)
	api.Post(registerRoute["registrationUpdate"], authz.RequirePermission(permissions.UpdateRegister), Update)
	api.Post(registerRoute["registrationDelete"], authz.RequirePermission(permissions.DeleteRegister), Delete)
	api.Post(registerRoute["registrationList"], authz.RequirePermission(permissions.ListRegisters), List)
	api.Post(registerRoute["registrationCancel"], authz.RequirePermission(permissions.CancelRegister), Cancel)
	api.Post(registerRoute["registrationListStudent"], authz.RequirePermission(permissions.ListStudentRegisters), ListStudent)
	api.Post(registerRoute["registrationListOffering"], authz.RequirePermission(permissions.ListOfferingRegisters), ListOffering)
	api.Post(registerRoute["registrationListMyClasses"], authz.RequirePermission(permissions.ListClasses), ClassesStudent)
	return registerRoute
}
