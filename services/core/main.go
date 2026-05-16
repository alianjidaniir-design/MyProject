package main

import (
	"MyProject/services/core/route"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	route.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
