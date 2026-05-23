package route

import (
	. "MyProject/controllers/student"

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
}

func SetupUserRoute(app *fiber.App) map[string]string {
	app.Post(studentRoute["studentCreate"], Create)
	app.Post(studentRoute["studentList"], List)
	app.Post(studentRoute["studentGet"], Get)
	app.Post(studentRoute["studentUpdate"], Update)
	app.Post(studentRoute["studentDelete"], Delete)
	app.Post(studentRoute["studentDelete2"], SoftDelete)
	app.Post(studentRoute["studentLogin"], Login)
	app.Post(studentRoute["studentRefresh"], Refresh)
	return studentRoute
}
