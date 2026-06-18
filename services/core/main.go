package main

import (
	"MyProject/services/core/route"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		return nil
	})
	route.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
