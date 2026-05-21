package authz

import (
	"MyProject/models/rolePermission/dataSources/mySQLDS"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var jwtSecretKey = []byte("your-super-secret-key") // <--- **این را با کلید واقعی خود جایگزین کنید**

func RequirePermission(permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID := GetRoleID(c)
		if roleID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token header is missing",
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
		joinQuery := fmt.Sprintf("SELECT * FROM rolepermissions JOIN permissions ON permissions.id = role_permissions.permission_id WHERE role_id = ? AND permission.name = ?")
		err = DB.QueryRow(joinQuery, roleID, permissionName).Scan(&count)
		if err != nil {
			return errors.New("123,d,sd")
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
