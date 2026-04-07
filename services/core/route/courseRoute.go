package route

import (
	. "MyProject/controllers/course"

	"github.com/gofiber/fiber/v2"
)

var routeCourse = map[string]string{
	"courseCreate": "course/create",
}

func SetupCourseRoutes(app *fiber.App) map[string]string {
	app.Post(routeCourse["courseCreate"], Create)
	return routeCourse
}
