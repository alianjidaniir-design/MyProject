package route

import (
	. "MyProject/controllers/category"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var categoryRoute = map[string]string{
	"categoryCreate": "category/create",
	"categoryDelete": "category/delete",
	"categoryGet":    "category/get",
	"categoryList":   "category/list",
}

func SetupCategoryRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(categoryRoute["categoryCreate"], authz.RequirePermission(permissions.CreateCategory), Create)
	api.Post(categoryRoute["categoryDelete"], authz.RequirePermission(permissions.DeleteCategory), Delete)
	api.Post(categoryRoute["categoryGet"], authz.RequirePermission(permissions.ViewCategory), Get)
	api.Post(categoryRoute["categoryList"], authz.RequirePermission(permissions.ListCategory), List)
	return categoryRoute
}
