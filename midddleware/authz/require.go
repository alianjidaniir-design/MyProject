package authz

import (
	"MyProject/models/student/dataModel"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("your-super-secret-key") // <--- **این را با کلید واقعی خود جایگزین کنید**

func RequirePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleHeader := c.Get("Authorization")
		if roleHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token header is missing",
			})
		}

		tokenStr := strings.TrimPrefix(roleHeader, "Bearer")

		claim := &dataModel.StudentRole{}

		token, err := jwt.ParseWithClaims(tokenStr, claim, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token is invalid",
			})
		}

		ctx := context.WithValue(context.Background(), "student_id", claim.StudentID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
