package route

import (
	. "MyProject/controllers/teacher"

	"github.com/gofiber/fiber/v2"
)

var teacherRoute = map[string]string{
	"TeacherCreate":     "teacher/create",
	"TeacherList":       "teacher/list",
	"TeacherDetail":     "teacher/detail",
	"TeacherDelete":     "teacher/delete",
	"TeacherSoftDelete": "teacher/soft_delete",
}

func SetupTeacherRoute(app *fiber.App) map[string]string {
	app.Post(teacherRoute["TeacherCreate"], Create)
	app.Post(teacherRoute["TeacherList"], List)
	app.Post(teacherRoute["TeacherDetail"], Get)
	app.Post(teacherRoute["TeacherDelete"], Delete)
	app.Post(teacherRoute["TeacherSoftDelete"], SoftDelete)

	return teacherRoute
}
