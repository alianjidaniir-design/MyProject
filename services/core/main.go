package main

import (
	"MyProject/services/core/route"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes := route.SetupRoutes(app)
	fmt.Println("very thanks Ali", routes)
	if err := app.Listen(":8080"); err != nil {
		fmt.Println(err)
	}
}
