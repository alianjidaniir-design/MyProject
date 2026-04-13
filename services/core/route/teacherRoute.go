package route

import (
	. "MyProject/controllers/teacher"

	"github.com/gofiber/fiber/v2"
)

var teacherRoute = map[string]string{
	"TeacherCreate": "teacher/create",
}

func SetupTeacherRoute(app *fiber.App) map[string]string {
	app.Post(teacherRoute["TeacherCreate"], Create)
	return teacherRoute
}
