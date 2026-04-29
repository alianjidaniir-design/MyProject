package route

import (
	. "MyProject/controllers/tuition"

	"github.com/gofiber/fiber/v2"
)

var tuitionRoute = map[string]string{
	"TuitionCreate": "tuition/create",
	"TuitionUpdate": "tuition/update",
}

func SetupTuitionRoute(app *fiber.App) map[string]string {
	app.Post(tuitionRoute["TuitionCreate"], Create)
	app.Post(tuitionRoute["TuitionUpdate"], Update)
	return tuitionRoute
}
