package authz

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (a *AuthzMiddleWare) RequirePermission(p string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleHeader := c.Get("role")
		if roleHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "role header is missing",
			})
		}

		roleParse, err := a.ParseRole(roleHeader)
		if err != nil {
			return err
		}

		check, err := a.HasPermissionByTID(roleParse, p)
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
