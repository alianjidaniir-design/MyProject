package route

import (
	. "MyProject/controllers/author"

	"github.com/gofiber/fiber/v2"
)

var authorRoute = map[string]string{
	"authorCreate": "author/create",
}

func SetupAuthorRoute(app *fiber.App) map[string]string {
	app.Post(authorRoute["authorCreate"], Create)
	return bookRoute
}
