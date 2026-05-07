package dataModel

type Program struct {
	Row         int64  `json:"row"`
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
