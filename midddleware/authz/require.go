package authz

import (
	"MyProject/statics/constants/permissions"
	"MyProject/statics/constants/status"

	"github.com/gofiber/fiber/v2"
)

func requirePermission(permissions permissions.Permissions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleHeader := c.Get("X-Role")
		if roleHeader == "" {
			return c.Status(status.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing role",
				"code":    "AUTH_01",
			})
		}
		role, err := ParseRole(roleHeader)
		if err != nil {
			return err
		}
		if !HasPermissionByTID(role, permissions) {
			return c.Status(status.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
				"code":    "AUTH_02",
			})
		}
		return c.Next()
	}
}
