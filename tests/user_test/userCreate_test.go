package user_test

import (
	"MyProject/services/core/route"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCreateUser_WithAuthz(t *testing.T) {
	app := fiber.New()

	// middleware تستی برای شبیه‌سازی کاربر مجاز
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_id", "admin")
		c.Locals("permissions", map[string]bool{
			"user/create": true,
		})
		return c.Next()
	})

	route.SetupRoutes(app)

	payload := map[string]any{
		"body": map[string]any{
			"code":   "313",
			"name":   "seyed",
			"family": "saeed",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload failed: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}
