package route

import (
	. "MyProject/controllers/category"

	"github.com/gofiber/fiber/v2"
)

var categoryRoute = map[string]string{
	"categoryCreate": "category/create",
	"categoryDelete": "category/delete",
}

func SetupCategoryRoute(app *fiber.App) map[string]string {
	app.Post(categoryRoute["categoryCreate"], Create)
	app.Post(categoryRoute["categoryDelete"], Delete)
	return categoryRoute
}
