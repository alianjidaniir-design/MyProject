package route

import (
	. "MyProject/controllers/teacher"
	"MyProject/midddleware/authz"

	"github.com/gofiber/fiber/v2"
)

var teacherRoute = map[string]string{
	"TeacherCreate":     "teacher/create",
	"TeacherList":       "teacher/list",
	"TeacherDetail":     "teacher/detail",
	"TeacherDelete":     "teacher/delete",
	"TeacherSoftDelete": "teacher/soft_delete",
	"TeacherUpdate":     "teacher/update",
	"TeacherLogin":      "teacher/login",
	"TeacherRefresh":    "teacher/refresh",
	"TeacherLogout":     "teacher/logout",
}

func SetupTeacherRoute(app *fiber.App) map[string]string {
	app.Post(teacherRoute["TeacherCreate"], Create)
	app.Post(teacherRoute["TeacherList"], List)
	app.Post(teacherRoute["TeacherDetail"], Get)
	app.Post(teacherRoute["TeacherDelete"], Delete)
	app.Post(teacherRoute["TeacherSoftDelete"], SoftDelete)
	app.Post(teacherRoute["TeacherUpdate"], Update)
	app.Post(teacherRoute["TeacherLogin"], Login)
	app.Post(teacherRoute["TeacherRefresh"], Refresh)
	app.Post(teacherRoute["TeacherLogout"], authz.AuthMiddleware(), Logout)

	return teacherRoute
}
