package user_test

import (
	"MyProject/services/core/route"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp(withPermission bool, roleName string) *fiber.App {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {

		if withPermission {
			headerRoleName := c.Get("role")
			if headerRoleName == "" {
				headerRoleName = roleName
			}
			if headerRoleName == "Admin" {
				c.Locals("role_id", 4)
				c.Locals("permissions", map[string]bool{
					"user/create": true,
					"user/read":   true,
				})
			} else if headerRoleName == "Student" {
				c.Locals("role_id", 3)
				c.Locals("permissions", map[string]bool{
					"user/read": true,
				})
			} else {
				// اگر نقش ناشناخته بود، اجازه دسترسی نده
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid role",
				})
			}
		}

		return c.Next()
	})

	route.SetupRoutes(app)
	return app
}

func TestCreateUser_WithAuthz(t *testing.T) {
	app := setupTestApp(true, "Admin")

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

	req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("role", "Student")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 200 or 201, got %d", resp.StatusCode)
	}
}

func TestCreateUser_WithoutAuthz(t *testing.T) {
	app := setupTestApp(false, "Student")

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

	req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("role", "Student")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401 or 403, got %d", resp.StatusCode)
	}
}

func TestCreateUser_MissingAuthHeader(t *testing.T) {
	app := setupTestApp(true, "") // با احراز هویت فعال، ولی انتظار نقش از هدر

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

	req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	// بررسی status code (باید 401 باشد چون هدر نقش نیست)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

// این تست برای حالتی است که authentication کلاً غیرفعال است
func TestCreateUser_NoAuthMiddleware(t *testing.T) {
	app := setupTestApp(false, "") // احراز هویت غیرفعال است

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

	req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	// چون احراز هویتی نیست، باید موفقیت‌آمیز باشد (status 200 یا 201)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 200 or 201 when auth is disabled, got %d", resp.StatusCode)
	}
}
