package route

import (
	. "MyProject/controllers/publisher"

	"github.com/gofiber/fiber/v2"
)

var PublisherRoute = map[string]string{
	"PublisherCreate": "publisher/create",
	"PublisherDetail": "publisher/detail",
}

func SetupPublisherRoute(app *fiber.App) map[string]string {
	app.Post(PublisherRoute["PublisherCreate"], Create)
	app.Post(PublisherRoute["PublisherDetail"], Detail)
	return PublisherRoute
}
