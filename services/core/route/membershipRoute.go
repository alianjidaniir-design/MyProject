package route

import (
	. "MyProject/controllers/membership"

	"github.com/gofiber/fiber/v2"
)

var membershipRoute = map[string]string{
	"membershipCreate": "membership/create",
}

func SetupMembershipRoute(app *fiber.App) map[string]string {
	app.Post(membershipRoute["membershipCreate"], Create)
	return membershipRoute
}
