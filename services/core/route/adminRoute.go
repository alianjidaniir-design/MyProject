package route

import (
	. "MyProject/controllers/admin"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

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
	api.Post(adminRoute["adminCreate"], authz.RequirePermission(permissions.CreateAdmin), Create)
	api.Post(adminRoute["adminLogout"], Logout)
	return adminRoute
}
