package route

import (
	. "MyProject/controllers/category"

	"github.com/gofiber/fiber/v2"
)

var categoryRoute = map[string]string{
	"categoryCreate": "category/create",
	"categoryDelete": "category/delete",
	"categoryGet":    "category/get",
	"categoryList":   "category/list",
}

func SetupCategoryRoute(app *fiber.App) map[string]string {
	app.Post(categoryRoute["categoryCreate"], Create)
	app.Post(categoryRoute["categoryDelete"], Delete)
	app.Post(categoryRoute["categoryGet"], Get)
	app.Post(categoryRoute["categoryList"], List)
	return categoryRoute
}
