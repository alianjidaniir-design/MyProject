package route

import (
	. "MyProject/controllers/course"

	"github.com/gofiber/fiber/v2"
)

var routeCourse = map[string]string{
	"courseCreate": "course/create",
	"courseList":   "course/list",
	"courseDetail": "course/detail",
	"courseUpdate": "course/update",
	"courseDelete": "course/delete",
}

func SetupCourseRoutes(app *fiber.App) map[string]string {
	app.Post(routeCourse["courseCreate"], Create)
	app.Post(routeCourse["courseList"], List)
	app.Post(routeCourse["courseDetail"], Get)
	app.Post(routeCourse["courseUpdate"], Update)
	app.Post(routeCourse["courseDelete"], Delete)

	return routeCourse
}
