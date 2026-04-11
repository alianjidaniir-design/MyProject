package route

import (
	. "MyProject/controllers/enrollment"

	"github.com/gofiber/fiber/v2"
)

var enrollRoute = map[string]string{
	"enrollmentcreate": "enrollment/create",
	"enrollmentcancel": "enrollment/cancel",
	"enrollmentlist":   "enrollment/list",
}

func SetupEnrollmentRoute(app *fiber.App) map[string]string {
	app.Post(enrollRoute["enrollmentcreate"], Create)
	app.Post(enrollRoute["enrollmentcancel"], Cancel)
	app.Post(enrollRoute["enrollmentlist"], List)

	return enrollRoute
}
