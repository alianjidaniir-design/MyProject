package route

import (
	. "MyProject/controllers/program"

	"github.com/gofiber/fiber/v2"
)

var programRoute = map[string]string{
	"ProgramCreate": "program/create",
	"ProgramGet":    "program/get",
}

func SetupProgramRoute(app *fiber.App) map[string]string {
	app.Post(programRoute["ProgramCreate"], Create)
	app.Post(programRoute["ProgramGet"], Get)
	return programRoute
}
