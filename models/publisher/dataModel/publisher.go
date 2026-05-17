package dataModel

type Publisher struct {
	ID      int64  `json:"ID"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
