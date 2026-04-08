package user

import (
	"MyProject/services/core/route"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestDeleteUser(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)
	payload := map[string]any{
		"body": map[string]any{
			"ID": 126,
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal paylpad failed : %v", err)
	}

	req, err := http.NewRequest("POST", "/user/delete", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("create request failed : %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request failed : %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("test request failed with status code %d", resp.StatusCode)
	}

}
