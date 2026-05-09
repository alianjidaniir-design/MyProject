package route

import (
	. "MyProject/controllers/program"

	"github.com/gofiber/fiber/v2"
)

var programRoute = map[string]string{
	"ProgramCreate": "program/create",
	"ProgramGet":    "program/get",
	"ProgramDelete": "program/delete",
	"ProgramList":   "program/list",
}

func SetupProgramRoute(app *fiber.App) map[string]string {
	app.Post(programRoute["ProgramCreate"], Create)
	app.Post(programRoute["ProgramGet"], Get)
	app.Post(programRoute["ProgramDelete"], Delete)
	app.Post(programRoute["ProgramList"], List)
	return programRoute
}
