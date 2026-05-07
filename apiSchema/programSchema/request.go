package programSchema

type CreateProgramRequest struct {
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetDetailProgramRequest struct {
	Row int64 `json:"row"`
}
