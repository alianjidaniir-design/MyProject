package route

import (
	. "MyProject/controllers/author"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var authorRoute = map[string]string{
	"authorCreate": "author/create",
	"authorGet":    "author/get",
	"authorDelete": "author/delete",
	"authorList":   "author/list",
}

func SetupAuthorRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(authorRoute["authorCreate"], authz.RequirePermission(permissions.CreateAuthor), Create)
	api.Post(authorRoute["authorGet"], authz.RequirePermission(permissions.ViewAuthor), Get)
	api.Post(authorRoute["authorDelete"], authz.RequirePermission(permissions.DeleteAuthor), Delete)
	api.Post(authorRoute["authorList"], authz.RequirePermission(permissions.ListAuthors), List)
	return authorRoute
}
