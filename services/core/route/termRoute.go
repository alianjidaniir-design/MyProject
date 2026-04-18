package route

import (
	. "MyProject/controllers/term"
	"github.com/gofiber/fiber/v2"
)

var termRoute = map[string]string{
	"termCreate": "term/create",
	"termList":   "term/list",
}

func SetupTermRoute(app *fiber.App) map[string]string {
	app.Post(termRoute["termCreate"], Create)
	app.Post(termRoute["termList"], List)
	return termRoute
}
