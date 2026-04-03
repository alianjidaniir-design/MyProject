package route

import (
	. "MyProject/controllers/user"

	"github.com/gofiber/fiber/v2"
)

var userRoute = map[string]string{
	"userCreate": "/user/create",
	"userList":   "/user/list",
}

func SetupUserRoute(app *fiber.App) map[string]string {
	app.Post(userRoute["userCreate"], Create)
	app.Post(userRoute["userList"], List)
	return userRoute
}
