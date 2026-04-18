package route

import (
	. "MyProject/controllers/department"

	"github.com/gofiber/fiber/v2"
)

var routeDepartment = map[string]string{
	"departmentCreate": "department/create",
}

func SetupDepartmentRoutes(app *fiber.App) map[string]string {
	app.Post(routeDepartment["departmentCreate"], Create)
	return routeCourse
}
