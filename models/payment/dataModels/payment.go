package dataModels

import (
	"time"
)

type Payment struct {
	ID                int64      `json:"id"`
	TuitionRow        int64      `json:"tuition_row"`
	PaymentType       string     `json:"payment_type"`
	NumberInstallment NullString `json:"number_installment"`
	InstallmentTotal  NullInt64  `json:"installment_total"`
	InstallmentAmount NullInt64  `json:"installment_amount"`
	CashAmount        NullInt64  `json:"cash_amount"`
	Bank              string     `json:"bank"`
	Operation         bool       `json:"operation"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         time.Time  `json:"deleted_at"`
}

type NullInt64 struct {
	Val int64 `json:"value,omitempty"`
}

type NullString struct {
	Val string `json:"value,omitempty"`
}
