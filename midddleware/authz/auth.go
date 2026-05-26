package authz

import (
	"MyProject/models/token/dataModel"
	"MyProject/statics/configs"
	"errors"
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
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")

			}
			return jwtSecret, nil
		})
		if err != nil {
			// اگر خطا مربوط به اتمام زمان توکن باشد
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"code":    fiber.StatusUnauthorized,
					"message": "token has expired",
				})
			}
			// برای خطاهای دیگر توکن
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "invalid token: " + err.Error(),
			})
		}

		// اگر توکن معتبر نباشد (این شرط بعد از بررسی err چک می‌شود)
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token validation failed",
			})
		}

		if claim == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "failed to get claims",
			})
		}

		userIDInt := claim.UserID
		roleIDInt := claim.RoleID
		scope := claim.Scope
		if scope != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "token is invalid",
			})
		}

		jti := claim.ID
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

		c.Locals("user_id", userIDInt)
		c.Locals("role_id", roleIDInt)

		c.Locals("jti", jti)
		c.Locals("exp", claim.ExpiresAt.Time)

		return c.Next()
	}
}

// ای پی ای محافظت شده
func GetUserID(c *fiber.Ctx) int64 {
	if k, ok := c.Locals("user_id").(int64); ok {
		return k
	}
	return 0
}

// ای پی ای رول پرمیشن
func GetRoleID(c *fiber.Ctx) int64 {
	if k, ok := c.Locals("role_id").(int64); ok {
		return k
	}
	return 0
}
