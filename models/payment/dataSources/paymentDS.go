package dataSources

import (
	"MyProject/apiSchema/paymentSchema"
	"MyProject/models/payment/dataModels"
	"context"
)

type PaymentDS interface {
	CreatePayment(ctx context.Context, req paymentSchema.ConfirmationSchema) (res dataModels.Payment, err error)
	DeletePayment(ctx context.Context, req paymentSchema.DeleteInformation) (res dataModels.Payment, err error)
}
