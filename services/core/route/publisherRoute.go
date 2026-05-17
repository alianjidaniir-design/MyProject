package route

import (
	. "MyProject/controllers/publisher"

	"github.com/gofiber/fiber/v2"
)

var PublisherRoute = map[string]string{
	"PublisherCreate": "publisher/create",
	"PublisherDetail": "publisher/detail",
	"PublisherDelete": "publisher/delete",
	"PublisherList":   "publisher/list",
}

func SetupPublisherRoute(app *fiber.App) map[string]string {
	app.Post(PublisherRoute["PublisherCreate"], Create)
	app.Post(PublisherRoute["PublisherDetail"], Detail)
	app.Post(PublisherRoute["PublisherDelete"], Delete)
	app.Post(PublisherRoute["PublisherList"], List)
	return PublisherRoute
}
