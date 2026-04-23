package route

import (
	. "MyProject/controllers/profile"

	"github.com/gofiber/fiber/v2"
)

var profileRoute = map[string]string{
	"ProfileCreate": "profile/create",
}

func SetupProfileRoute(app *fiber.App) map[string]string {
	app.Post(profileRoute["ProfileCreate"], Create)
	
	return profileRoute
}
