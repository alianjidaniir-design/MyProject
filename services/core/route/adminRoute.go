package route

import (
	. "MyProject/controllers/admin"
	"MyProject/midddleware/authz"

	"github.com/gofiber/fiber/v2"
)

var adminRoute = map[string]string{
	"adminCreate":  "admin/create",
	"adminLogin":   "admin/login",
	"adminRefresh": "admin/refresh",
	"adminLogout":  "admin/logout",
}

func SetupAdminRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	app.Post(adminRoute["adminLogin"], Login)
	app.Post(adminRoute["adminRefresh"], Refresh)
	api.Post(adminRoute["adminCreate"], Create)
	api.Post(adminRoute["adminLogout"], Logout)
	api.Post(adminRoute["adminLogout"], Logout)
	return adminRoute
}
