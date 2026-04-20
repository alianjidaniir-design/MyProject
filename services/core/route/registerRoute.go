package route

import (
	. "MyProject/controllers/registration"

	"github.com/gofiber/fiber/v2"
)

var registerRoute = map[string]string{
	"registrationCreate": "registration/create",
}

func SetupRegistrationRoute(app *fiber.App) map[string]string {
	app.Post(registerRoute["registrationCreate"], Create)
	return registerRoute
}
