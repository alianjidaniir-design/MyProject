package paymentSchema

import "MyProject/models/payment/dataModels"

type DetailPaymentSchema struct {
	Detail dataModels.Payment `json:"detail"`
}

type DetailChangePaymentSchema struct {
	Detail  dataModels.Payment `json:"detail"`
	Massage string             `json:"massage"`
}
