package route

import (
	. "MyProject/controllers/department"

	"github.com/gofiber/fiber/v2"
)

var routeDepartment = map[string]string{
	"departmentCreate": "department/create",
	"departmentUpdate": "department/update",
	"departmentList":   "department/list",
	"departmentDelete": "department/delete",
}

func SetupDepartmentRoutes(app *fiber.App) map[string]string {
	app.Post(routeDepartment["departmentCreate"], Create)
	app.Post(routeDepartment["departmentUpdate"], Update)
	app.Post(routeDepartment["departmentList"], List)
	app.Post(routeDepartment["departmentDelete"], Delete)
	return routeCourse
}
