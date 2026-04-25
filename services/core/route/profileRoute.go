package route

import (
	. "MyProject/controllers/profile"

	"github.com/gofiber/fiber/v2"
)

var profileRoute = map[string]string{
	"ProfileCreate":            "profile/create",
	"ProfileListScoresStudent": "profile/scorse/students",
	"ProfileSummery":           "profile/summery",
	"ProfileDetail":            "profile/detail",
	"ProfileDelete":            "profile/delete",
}

func SetupProfileRoute(app *fiber.App) map[string]string {
	app.Post(profileRoute["ProfileCreate"], Create)
	app.Post(profileRoute["ProfileListScoresStudent"], ListScoresStudents)
	app.Post(profileRoute["ProfileSummery"], ListSummeryStudents)
	app.Post(profileRoute["ProfileDetail"], Get)
	app.Post(profileRoute["ProfileDelete"], Delete)
	return profileRoute
}
