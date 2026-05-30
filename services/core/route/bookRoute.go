package route

import (
	. "MyProject/controllers/book"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var bookRoute = map[string]string{
	"bookCreate": "/book/create",
	"bookDelete": "/book/delete",
	"bookDetail": "/book/detail",
	"bookList":   "/book/list",
}

func SetupBookRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(bookRoute["bookCreate"], Create)
	api.Post(bookRoute["bookDelete"], Delete)
	api.Post(bookRoute["bookDetail"], Get)
	api.Post(bookRoute["bookList"], authz.RequirePermission(permissions.ListBooks), List)
	return bookRoute
}
