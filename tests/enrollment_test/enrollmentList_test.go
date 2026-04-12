package enrollment_test

import (
	"MyProject/services/core/route"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestListEnrollment(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)
	payload := map[string]any{
		"body": map[string]any{
			"page":     1,
			"pageSize": 43,
		},
	}

	marshal, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("error in marshal : %v", err)
	}
	req, err := http.NewRequest("POST", "/enrollment/list", bytes.NewBuffer(marshal))
	if err != nil {
		t.Fatalf("error in create request : %v", err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error in response : %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("error in response code : %v", resp.StatusCode)
	}

}
