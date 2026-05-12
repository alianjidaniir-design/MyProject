package InventorySchema

type CreateInventory struct {
	BookCode   int64  `json:"book_code"`
	NumberBook int    `json:"number_book"`
	Situation  string `json:"situation"`
}
