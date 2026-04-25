package route

import (
	. "MyProject/controllers/tuition"

	"github.com/gofiber/fiber/v2"
)

var toitionRoute = map[string]string{
	"TuitionCreate": "tuition/create",
}

func SetupTuitionRoute(app *fiber.App) map[string]string {
	app.Post(toitionRoute["TuitionCreate"], Create)
	return toitionRoute
}
