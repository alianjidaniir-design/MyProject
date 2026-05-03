package route

import (
	. "MyProject/controllers/tuition"

	"github.com/gofiber/fiber/v2"
)

var tuitionRoute = map[string]string{
	"TuitionCreate":       "tuition/create",
	"TuitionUpdate":       "tuition/update",
	"TuitionDelete":       "tuition/delete",
	"TuitionListStudents": "tuition/list/students",
	"TuitionList":         "tuition/list",
	"TuitionDetail":       "tuition/detail",
}

func SetupTuitionRoute(app *fiber.App) map[string]string {
	app.Post(tuitionRoute["TuitionCreate"], Create)
	app.Post(tuitionRoute["TuitionUpdate"], Update)
	app.Post(tuitionRoute["TuitionDelete"], Delete)
	app.Post(tuitionRoute["TuitionListStudents"], ListStudents)
	app.Post(tuitionRoute["TuitionList"], List)
	app.Post(tuitionRoute["TuitionDetail"], Get)
	return tuitionRoute
}
