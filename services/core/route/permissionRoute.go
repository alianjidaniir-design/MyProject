package route

import (
	. "MyProject/controllers/permission"

	"github.com/gofiber/fiber/v2"
)

var permissionRoute = map[string]string{
	"PermissionCreate": "permission/create",
	"PermissionDelete": "permission/delete",
	"PermissionDetail": "permission/detail",
	"PermissionList":   "permission/list",
}

func SetupPermissionRoute(app *fiber.App) map[string]string {
	app.Post(permissionRoute["PermissionCreate"], Create)
	app.Post(permissionRoute["PermissionDetail"], Get)
	app.Post(permissionRoute["PermissionDelete"], Delete)
	app.Post(permissionRoute["PermissionList"], List)
	return roleRoute
}
