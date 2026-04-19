package route

import (
	. "MyProject/controllers/offering"

	"github.com/gofiber/fiber/v2"
)

var offeringRoute = map[string]string{
	"offeringCreate": "offering/create",
	"offeringList":   "offering/list",
	"offeringDetail": "offering/detail",
}

func SetupOfferingRoute(app *fiber.App) map[string]string {
	app.Post(offeringRoute["offeringCreate"], Create)
	app.Post(offeringRoute["offeringList"], List)
	app.Post(offeringRoute["offeringDetail"], Get)
	return offeringRoute
}
