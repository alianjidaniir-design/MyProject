package route

import (
	. "MyProject/Controllers/subject"

	"github.com/gofiber/fiber/v2"
)

var subjectRoute = map[string]string{
	"SubjectCreate": "subject/create",
	"SubjectGet":    "subject/get",
	"SubjectDelete": "subject/delete",
	"SubjectList":   "subject/list",
}

func SetupSubjectRoute(app *fiber.App) map[string]string {
	app.Post(subjectRoute["SubjectCreate"], Create)
	app.Post(subjectRoute["SubjectGet"], Get)
	app.Post(subjectRoute["SubjectDelete"], Delete)
	app.Post(subjectRoute["SubjectList"], List)
	return subjectRoute
}
