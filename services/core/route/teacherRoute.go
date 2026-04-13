package route

import (
	. "MyProject/controllers/teacher"

	"github.com/gofiber/fiber/v2"
)

var teacherRoute = map[string]string{
	"TeacherCreate": "teacher/create",
	"TeacherList":   "teacher/list",
}

func SetupTeacherRoute(app *fiber.App) map[string]string {
	app.Post(teacherRoute["TeacherCreate"], Create)
	app.Post(teacherRoute["TeacherList"], List)
	return teacherRoute
}
