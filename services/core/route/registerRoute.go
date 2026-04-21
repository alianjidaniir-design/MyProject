package route

import (
	. "MyProject/controllers/registration"

	"github.com/gofiber/fiber/v2"
)

var registerRoute = map[string]string{
	"registrationCreate": "registration/create",
	"registrationGet":    "registration/get",
	"registrationUpdate": "registration/update",
	"registrationDelete": "registration/delete",
	"registrationList":   "registration/list",
}

func SetupRegistrationRoute(app *fiber.App) map[string]string {
	app.Post(registerRoute["registrationCreate"], Create)
	app.Post(registerRoute["registrationGet"], Get)
	app.Post(registerRoute["registrationUpdate"], Update)
	app.Post(registerRoute["registrationDelete"], Delete)
	app.Post(registerRoute["registrationList"], List)
	return registerRoute
}
