package paymentSchema

type ConfirmationSchema struct {
	TuitionRow        int64  `json:"tuition_row"`
	PaymentType       string `json:"payment_type"`
	NumberInstallment string `json:"number_installment,omitempty"`
	Bank              string `json:"bank"`
	Operation         bool   `json:"operation"`
}

type DeleteInformation struct {
	ID int64 `json:"id"`
}

type GetInformation struct {
	ID int64 `json:"id"`
}

type Filter struct {
	PaymentType string `json:"payment_type"`
	Bank        string `json:"bank"`
	Operation   bool   `json:"operation"`
}

type ListPayment struct {
	Filter *Filter `json:"filter"`
}
