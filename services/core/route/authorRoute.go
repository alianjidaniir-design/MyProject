package route

import (
	. "MyProject/controllers/author"

	"github.com/gofiber/fiber/v2"
)

var authorRoute = map[string]string{
	"authorCreate": "author/create",
	"authorGet":    "author/get",
	"authorDelete": "author/delete",
	"authorList":   "author/list",
}

func SetupAuthorRoute(app *fiber.App) map[string]string {
	app.Post(authorRoute["authorCreate"], Create)
	app.Post(authorRoute["authorGet"], Get)
	app.Post(authorRoute["authorDelete"], Delete)
	app.Post(authorRoute["authorList"], List)
	return bookRoute
}
