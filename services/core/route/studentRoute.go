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
}

func SetupUserRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(studentRoute["studentCreate"], authz.RequirePermission(permissions.CreateUser), Create)
	api.Post(studentRoute["studentList"], List)
	api.Post(studentRoute["studentGet"], authz.RequirePermission(permissions.ViewDetailStudent), Get)
	api.Post(studentRoute["studentUpdate"], Update)
	api.Post(studentRoute["studentDelete"], Delete)
	api.Post(studentRoute["studentDelete2"], SoftDelete)
	app.Post(studentRoute["studentLogin"], Login)
	app.Post(studentRoute["studentRefresh"], Refresh)
	api.Post(studentRoute["studentLogout"], authz.AuthMiddleware(), Logout)
	return studentRoute
}
