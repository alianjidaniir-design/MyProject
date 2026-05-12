package route

import (
	. "MyProject/controllers/role"

	"github.com/gofiber/fiber/v2"
)

var roleRoute = map[string]string{
	"RoleCreate": "role/create",
	"RoleDelete": "role/delete",
	"RoleDetail": "role/detail",
	"RoleList":   "role/list",
}

func SetupRoleRoute(app *fiber.App) map[string]string {
	app.Post(roleRoute["RoleCreate"], Create)
	app.Post(roleRoute["RoleDelete"], Delete)
	app.Post(roleRoute["RoleDetail"], Get)
	app.Post(roleRoute["RoleList"], List)
	return roleRoute
}
