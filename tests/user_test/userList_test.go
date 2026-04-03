package user

import (
	"MyProject/services/core/route"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func testListUser(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)
	createpayload := map[string]any{
		"body": map[string]any{
			"code":   "12345678",
			"name":   "John",
			"family": "razavi",
		},
	}
	createBody, err := json.Marshal(createpayload)
	if err != nil {
		t.Fatal("Error Marshal", err)
	}
	creatrReq, err := http.NewRequest("POST", "/user/create", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatal("Error Create Request", err)
	}
	creatrReq.Header.Set("Content-Type", "application/json")
	if _, err := app.Test(creatrReq); err != nil {
		t.Fatal("Error Test", err)
	}
	listreq, err := http.NewRequest("GET", "/user/list", nil)
	if err != nil {
		t.Fatal("Error List Request", err)
	}
	listRes, err := app.Test(listreq)
	if err != nil {
		t.Fatal("Error List Response", err)
	}
	if listRes.StatusCode != 200 {
		t.Fatal("Error List Response", listRes.StatusCode)
	}
}
