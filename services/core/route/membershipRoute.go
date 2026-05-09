package route

import (
	. "MyProject/controllers/membership"

	"github.com/gofiber/fiber/v2"
)

var membershipRoute = map[string]string{
	"membershipCreate":     "membership/create",
	"membershipDelete":     "membership/delete",
	"membershipUpdate":     "membership/update",
	"membershipDeactivate": "membership/deactivate",
	"membershipDetail":     "membership/detail",
}

func SetupMembershipRoute(app *fiber.App) map[string]string {
	app.Post(membershipRoute["membershipCreate"], Create)
	app.Post(membershipRoute["membershipDelete"], Delete)
	app.Post(membershipRoute["membershipUpdate"], Update)
	app.Post(membershipRoute["membershipDeactivate"], DeActive)
	app.Post(membershipRoute["membershipDetail"], Detail)
	return membershipRoute
}
