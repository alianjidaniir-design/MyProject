package user_test

import (
	. "MyProject/controllers/user"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Mock implementation for RoleDS
func setupTestApp() *fiber.App {
	app := fiber.New()

	// فقط داده‌های لازم برای handler را ست می‌کنیم
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_id", 4)
		c.Locals("permissions", map[string]bool{
			"user/create": true,
		})
		return c.Next()
	})

	// مستقیم همان handler را ثبت کن
	app.Post("/user/create", Create)

	return app
}

func TestCreateUser(t *testing.T) {
	app := setupTestApp()

	payload := map[string]any{
		"body": map[string]any{
			"code":   "313",
			"name":   "seyed",
			"family": "saeed",
		},
	}

	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("expected success, got %d", resp.StatusCode)
	}
}
