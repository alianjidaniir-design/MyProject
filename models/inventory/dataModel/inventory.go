package dataModel

type Inventory struct {
	ID         int64  `json:"id"`
	BookCode   int64  `json:"book_code"`
	NumberBook int    `json:"number_book"`
	Situation  string `json:"situation"`
}
