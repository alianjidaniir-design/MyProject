package authz

import (
	dataModel2 "MyProject/models/permission/dataModel"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RequirePermission(p dataModel2.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rileHeader := c.Get("role")
		if rileHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"massage": "this is role there is not",
			})
		}
		roleParse, err := ParseRole(rileHeader)
		if err != nil {
			return err
		}
		check, err := HasPermissionByTID(roleParse, p)
		if err != nil {
			return err
		}
		if !check {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    fiber.StatusForbidden,
				"message": fmt.Sprintf("role '%s' does not have permission '%s'", roleParse.Name, p),
			})
		}
		c.Locals("role", roleParse)
		c.Locals("role_id", roleParse.ID)
		return c.Next()

	}
}
