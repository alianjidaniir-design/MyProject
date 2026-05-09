package categorySchema

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type GetRowCategoryRequest struct {
	Row int64 `json:"row"`
}

type PaginationList struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
