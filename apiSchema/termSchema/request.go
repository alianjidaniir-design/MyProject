package termSchema

type CreateTerm struct {
	Term string `json:"term"`
	Year int    `json:"year"`
}
type ListTerm struct {
	PageIndex int `json:"page_index"`
	PageSize  int `json:"pageSize"`
}

type DeleteTerm struct {
	ID int `json:"id"`
}
