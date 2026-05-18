package route

import (
	. "MyProject/controllers/translator"

	"github.com/gofiber/fiber/v2"
)

var translatorRoute = map[string]string{
	"TranslatorCreate": "translator/create",
	"TranslatorGet":    "translator/get",
	"TranslatorDelete": "translator/delete",
	"TranslatorList":   "translator/list",
}

func SetupTranslatorRoute(app *fiber.App) map[string]string {
	app.Post(translatorRoute["TranslatorCreate"], Create)
	app.Post(translatorRoute["TranslatorGet"], Get)
	app.Post(translatorRoute["TranslatorDelete"], Delete)
	app.Post(translatorRoute["TranslatorList"], List)
	return translatorRoute
}
