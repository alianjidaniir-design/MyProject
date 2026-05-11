package route

import (
	. "MyProject/controllers/book"

	"github.com/gofiber/fiber/v2"
)

var bookRoute = map[string]string{
	"bookCreate": "/book/create",
	"bookDelete": "/book/delete",
	"bookDetail": "/book/detail",
	"bookList":   "/book/list",
}

func SetupBookRoute(app *fiber.App) map[string]string {
	app.Post(bookRoute["bookCreate"], Create)
	app.Post(bookRoute["bookDelete"], Delete)
	app.Post(bookRoute["bookDetail"], Get)
	app.Post(bookRoute["bookList"], List)
	return bookRoute
}
