package authz

import (
	"MyProject/models/student/dataModel"
	"MyProject/statics/configs"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(configs.AccessTokenSecret)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token header is missing",
			})
		}

		tokenStr := strings.Split(authHeader, " ")
		if len(tokenStr) != 2 || tokenStr[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "len token header is invalid",
			})
		}

		claim := &dataModel.AccessToken{}

		token, err := jwt.ParseWithClaims(tokenStr[1], claim, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token is invalid " + err.Error(),
			})
		}

		userIDInt := token.Claims.(*dataModel.AccessToken).UserID
		roleIDInt := token.Claims.(*dataModel.AccessToken).RoleID
		scope := token.Claims.(*dataModel.AccessToken).Scope
		if scope != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token is invalid",
			})
		}

		claims := token.Claims.(*dataModel.AccessToken)

		jti := claims.ID
		if jti == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "jti is empty",
			})
		}
		if claim.ExpiresAt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token expire invalid",
			})
		}

		c.Locals("userID", userIDInt)
		c.Locals("roleID", roleIDInt)

		c.Locals("jti", jti)
		c.Locals("exp", claims.ExpiresAt.Time)

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
