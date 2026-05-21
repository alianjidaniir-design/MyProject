package authz

import (
	"MyProject/models/student/dataModel"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("your-super-secret-key") // <--- **این را با کلید واقعی خود جایگزین کنید**

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token header is missing",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer")

		claim := &dataModel.AccessToken{}

		token, err := jwt.ParseWithClaims(tokenStr, claim, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token is invalid",
			})
		}

		userIDInt, ok1 := token.Claims.(jwt.MapClaims)["user_id"].(int64)
		roleIDInt, ok2 := token.Claims.(jwt.MapClaims)["role_id"].(int64)
		scope, ok3 := token.Claims.(jwt.MapClaims)["scope"].(string)
		if !ok1 || !ok2 || !ok3 || scope != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token is invalid",
			})
		}

		c.Locals("userID", userIDInt)
		c.Locals("roleID", roleIDInt)

		return c.Next()
	}
}

// ای پی ای محافظت شده
func GetUserID(c *fiber.Ctx) int64 {
	if k, ok := c.Locals("userID").(int64); ok {
		return k
	}
	return 0
}

// ای پی ای رول پرمیشن
func GetRoleID(c *fiber.Ctx) int64 {
	if k, ok := c.Locals("roleID").(int64); ok {
		return k
	}
	return 0
}
