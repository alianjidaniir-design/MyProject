package route

import (
	. "MyProject/controllers/user"

	"github.com/gofiber/fiber/v2"
)

var userRoute = map[string]string{
	"userCreate":  "/user/create",
	"userList":    "/user/list",
	"userGet":     "/user/get",
	"userUpdate":  "/user/update",
	"userDelete":  "/user/delete",
	"userDelete2": "/user/delete2",
}

func SetupUserRoute(app *fiber.App) map[string]string {
	app.Post(userRoute["userCreate"], Create)
	app.Post(userRoute["userList"], List)
	app.Post(userRoute["userGet"], Get)
	app.Post(userRoute["userUpdate"], Update)
	app.Post(userRoute["userDelete"], Delete)
	app.Post(userRoute["userDelete2"], SoftDelete)
	return userRoute
}
