package route

import (
	. "MyProject/controllers/profile"

	"github.com/gofiber/fiber/v2"
)

var profileRoute = map[string]string{
	"ProfileCreate":            "profile/create",
	"ProfileListScoresStudent": "profile/scorse/students",
}

func SetupProfileRoute(app *fiber.App) map[string]string {
	app.Post(profileRoute["ProfileCreate"], Create)
	app.Post(profileRoute["ProfileListScoresStudent"], ListScoresStudents)

	return profileRoute
}
