package route

import (
	. "MyProject/controllers/course"

	"github.com/gofiber/fiber/v2"
)

var routeCourse = map[string]string{
	"courseCreate": "course/create",
	"courseList":   "course/list",
	"courseDetail": "course/detail",
}

func SetupCourseRoutes(app *fiber.App) map[string]string {
	app.Post(routeCourse["courseCreate"], Create)
	app.Post(routeCourse["courseList"], List)
	app.Post(routeCourse["courseDetail"], Get)
	return routeCourse
}
