package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/paymentSchema"
	"MyProject/models/payment"
	"context"
)

type PaymentRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.ConfirmationSchema]) (res paymentSchema.DetailPaymentSchema, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.DeleteInformation]) (res paymentSchema.DetailChangePaymentSchema, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.GetInformation]) (res paymentSchema.DetailPaymentSchema, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.ListPayment]) (res paymentSchema.DetailListPaymentSchema, errStr string, code int, err error)
}

var PaymentRepo PaymentRepository = payment.GetRepo()
