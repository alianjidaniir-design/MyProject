package route

import (
	. "MyProject/controllers/enrollment"

	"github.com/gofiber/fiber/v2"
)

var enrollRoute = map[string]string{
	"enrollmentCreate":   "enrollment/create",
	"enrollmentCancel":   "enrollment/cancel",
	"enrollmentList":     "enrollment/list",
	"enrollmentCourses":  "enrollment/courses",
	"enrollmentStudents": "enrollment/students",
}

func SetupEnrollmentRoute(app *fiber.App) map[string]string {
	app.Post(enrollRoute["enrollmentCreate"], Create)
	app.Post(enrollRoute["enrollmentCancel"], Cancel)
	app.Post(enrollRoute["enrollmentList"], List)
	app.Post(enrollRoute["enrollmentCourses"], ListCourses)
	app.Post(enrollRoute["enrollmentStudents"], ListStudents)

	return enrollRoute
}
