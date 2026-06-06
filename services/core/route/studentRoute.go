package route

import (
	. "MyProject/controllers/student"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var studentRoute = map[string]string{
	"studentCreate":  "/student/create",
	"studentList":    "/student/list",
	"studentGet":     "/student/get",
	"studentUpdate":  "/student/update",
	"studentDelete":  "/student/delete",
	"studentDelete2": "/student/delete2",
	"studentLogin":   "/student/login",
	"studentRefresh": "/student/refresh",
	"studentLogout":  "/student/logout",
	"studentMyInfo":  "/student/myinfo",
}

func SetupUserRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(studentRoute["studentCreate"], authz.RequirePermission(permissions.CreateStudent), Create)
	api.Post(studentRoute["studentList"], authz.RequirePermission(permissions.ListStudent), List)
	api.Post(studentRoute["studentGet"], authz.RequirePermission(permissions.ViewDetailStudent), Get)
	api.Post(studentRoute["studentUpdate"], authz.RequirePermission(permissions.UpdateStudent), Update)
	api.Post(studentRoute["studentDelete"], authz.RequirePermission(permissions.DeleteStudent), Delete)
	api.Post(studentRoute["studentDelete2"], authz.RequirePermission(permissions.DeleteStudent), SoftDelete)
	app.Post(studentRoute["studentLogin"], Login)
	app.Post(studentRoute["studentRefresh"], Refresh)
	api.Post(studentRoute["studentLogout"], Logout)
	api.Post(studentRoute["studentMyInfo"], authz.RequirePermission(permissions.MyViewStudent), MyInfo)
	return studentRoute
}
