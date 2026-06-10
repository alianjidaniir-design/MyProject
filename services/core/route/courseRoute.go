package route

import (
	. "MyProject/controllers/course"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var routeCourse = map[string]string{
	"courseCreate":     "course/create",
	"courseList":       "course/list",
	"courseDetail":     "course/detail",
	"courseUpdate":     "course/update",
	"courseDelete":     "course/delete",
	"courseSoftDelete": "course/soft_delete",
	"courseDepartment": "course/list/department",
}

func SetupCourseRoutes(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(routeCourse["courseCreate"], authz.RequirePermission(permissions.CreateCourse), Create)
	api.Post(routeCourse["courseList"], authz.RequirePermission(permissions.ListCourses), List)
	api.Post(routeCourse["courseDetail"], authz.RequirePermission(permissions.ViewCourse), Get)
	api.Post(routeCourse["courseUpdate"], authz.RequirePermission(permissions.UpdateCourse), Update)
	api.Post(routeCourse["courseDelete"], authz.RequirePermission(permissions.DeleteCourse), Delete)
	api.Post(routeCourse["courseSoftDelete"], authz.RequirePermission(permissions.SoftDeleteCourse), SoftDelete)
	api.Post(routeCourse["courseDepartment"], authz.RequirePermission(permissions.ListDepartmentCourse), ListDepartment)

	return routeCourse
}
