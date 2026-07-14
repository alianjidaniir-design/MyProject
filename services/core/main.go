package main

import (
	"MyProject/apiSchema/adminSchema"
	"MyProject/services/core/route"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/salam", func(c *fiber.Ctx) error {
		a, _ := json.Marshal(adminSchema.DetailAdminSchema{})
		c.Write(a)
		return nil
	})
	route.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
