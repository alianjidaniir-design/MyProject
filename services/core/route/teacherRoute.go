package route

import (
	. "MyProject/controllers/teacher"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var teacherRoute = map[string]string{
	"TeacherCreate":     "teacher/create",
	"TeacherList":       "teacher/list",
	"TeacherDetail":     "teacher/detail",
	"TeacherSoftDelete": "teacher/soft_delete",
	"TeacherUpdate":     "teacher/update",
	"TeacherLogin":      "teacher/login",
	"TeacherRefresh":    "teacher/refresh",
	"TeacherLogout":     "teacher/logout",
	"TeacherMyInfo":     "teacher/myinfo",
}

func SetupTeacherRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(teacherRoute["TeacherCreate"], authz.RequirePermission(permissions.CreateTeacher), Create)
	api.Post(teacherRoute["TeacherList"], authz.RequirePermission(permissions.ListTeachers), List)
	api.Post(teacherRoute["TeacherDetail"], authz.RequirePermission(permissions.ViewDetailTeacher), Get)
	api.Post(teacherRoute["TeacherSoftDelete"], authz.RequirePermission(permissions.SoftDeleteTeacher), SoftDelete)
	api.Post(teacherRoute["TeacherUpdate"], authz.RequirePermission(permissions.UpdateTeacher), Update)
	app.Post(teacherRoute["TeacherLogin"], Login)
	app.Post(teacherRoute["TeacherRefresh"], Refresh)
	api.Post(teacherRoute["TeacherLogout"], Logout)
	api.Post(teacherRoute["TeacherMyInfo"], authz.RequirePermission(permissions.ViewDetailTeacher), MyInfo)

	return teacherRoute
}
