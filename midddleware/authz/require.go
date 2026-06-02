package authz

import (
	"MyProject/models/rolePermission/dataSources/mySQLDS"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RequirePermission(permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName := GetRoleName(c)
		if roleName == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "role not exist",
			})
		}

		cfg, err := mySQLDS.LoadConfig()
		if err != nil {
			return err
		}
		DB, err := mySQLDS.Open(cfg)
		if err != nil {
			return err
		}

		var count int64

		var countPer int64
		checkPermission := fmt.Sprintf("SELECT COUNT(*) FROM permissions WHERE name = ? ")
		err = DB.QueryRow(checkPermission, permissionName).Scan(&countPer)
		if err != nil {
			return err
		}
		if countPer == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "permission not exist",
			})
		}

		joinQuery := `SELECT COUNT(*) FROM rolepermissions 
        JOIN permissions ON permissions.id = rolepermissions.permission_id 
        JOIN roles ON roles.id = rolepermissions.role_id       
        WHERE roles.name = ? AND permissions.name = ?`
		err = DB.QueryRow(joinQuery, roleName, permissionName).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "permission denied",
			})
		}

		return c.Next()
	}
}
