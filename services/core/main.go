package main

import (
	"MyProject/services/core/route"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes := route.SetupRoutes(app)
	fmt.Println("Project of University", routes)
	if err := app.Listen(":3000"); err != nil {
		fmt.Println(err)
	}
}
