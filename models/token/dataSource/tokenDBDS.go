package dataSource

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Cookies(c *fiber.Ctx) (string, error) {
	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return "", errors.New("no refresh_token")
	}
	return refreshToken, nil
}

func BList(c *fiber.Ctx) (string, time.Time, error) {
	fmt.Println(c.Locals("jti"), c.Locals("exp"))
	jti, ok := c.Locals("jti").(string)
	exp, ok2 := c.Locals("exp").(time.Time)
	if !ok || !ok2 {
		return "", time.Time{}, errors.New("no jti")
	}
	return jti, exp, nil

}
