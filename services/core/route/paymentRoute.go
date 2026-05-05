package route

import (
	. "MyProject/controllers/payment"

	"github.com/gofiber/fiber/v2"
)

var paymentRoute = map[string]string{
	"paymentCreate": "payment/create",
	"paymentDelete": "payment/delete",
}

func SetupPaymentRoute(app *fiber.App) map[string]string {
	app.Post(paymentRoute["paymentCreate"], Create)
	app.Post(paymentRoute["paymentDelete"], Delete)
	return paymentRoute
}
