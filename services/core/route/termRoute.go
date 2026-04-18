package route

import (
	. "MyProject/controllers/term"
	"github.com/gofiber/fiber/v2"
)

var termRoute = map[string]string{
	"termCreate": "term/create",
}

func SetupTermRoute(app *fiber.App) map[string]string {
	app.Post(termRoute["termCreate"], Create)
	return termRoute
}
