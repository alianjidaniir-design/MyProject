package route

import (
	. "MyProject/controllers/rolePermission"

	"github.com/gofiber/fiber/v2"
)

var rolePermissionRoute = map[string]string{
	"rolePermissionCreate": "rolepermission/create",
}

func SetupRolePermissionRoute(app *fiber.App) map[string]string {
	app.Post(rolePermissionRoute["rolePermissionCreate"], Create)
	return rolePermissionRoute
}
