package termSchema

type CreateTerm struct {
	Term int `json:"term"`
	Year int `json:"year"`
}
type ListTerm struct {
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}

type DeleteTerm struct {
	ID int `json:"id"`
}
