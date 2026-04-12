package user

import (
	"MyProject/services/core/route"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestListUser(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)
	createPayload := map[string]any{
		"body": map[string]any{
			"page":   "12345678",
			"name":   "John",
			"family": "raze",
		},
	}
	listPayload := map[string]any{
		"body": map[string]any{
			"page":     1,
			"pageSize": 5,
		},
	}
	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatal("Error Marshal", err)
	}
	createReq, err := http.NewRequest("POST", "/user/create", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatal("Error Create Request", err)
	}
	createReq.Header.Set("Content-Type", "application/json")
	if _, err := app.Test(createReq); err != nil {
		t.Fatal("Error Test", err)
	}
	list, err := json.Marshal(listPayload)
	if err != nil {
		t.Fatal("Error Marshal", err)
	}
	listReq, err := http.NewRequest("POST", "/user/list", bytes.NewBuffer(list))
	if err != nil {
		t.Fatal("Error List Request", err)
	}
	listRes, err := app.Test(listReq)
	if err != nil {
		t.Fatal("Error List Response", err)
	}
	if listRes.StatusCode != 200 {
		t.Fatal("Error List Response", listRes.StatusCode)
	}
}
