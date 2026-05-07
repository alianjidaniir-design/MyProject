package categorySchema

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type GetRowCategoryRequest struct {
	Row int64 `json:"row"`
}
