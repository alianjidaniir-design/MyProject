package route

import (
	. "MyProject/controllers/offering"
	"MyProject/midddleware/authz"
	"MyProject/statics/constants/permissions"

	"github.com/gofiber/fiber/v2"
)

var offeringRoute = map[string]string{
	"offeringCreate":     "offering/create",
	"offeringList":       "offering/list",
	"offeringDetail":     "offering/detail",
	"offeringDeActivate": "offering/deactivate",
}

func SetupOfferingRoute(app *fiber.App) map[string]string {
	api := app.Group("/api", authz.AuthMiddleware())
	api.Post(offeringRoute["offeringCreate"], authz.RequirePermission(permissions.CreateOffering), Create)
	api.Post(offeringRoute["offeringList"], authz.RequirePermission(permissions.ListOfferings), List)
	api.Post(offeringRoute["offeringDetail"], authz.RequirePermission(permissions.ViewOffering), Get)
	api.Post(offeringRoute["offeringDeActivate"], authz.RequirePermission(permissions.DeActivateOffering), DeActive)
	return offeringRoute
}
