package profileSchema

type CreateScoresReq struct {
	RegistrationID int64 `json:"registration_id"`
	Score          int   `json:"score"`
}

type ListAllScoresReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetScoresReq struct {
	StudentID int64 `json:"student_id"`
}

type DeleteScoresReq struct {
	ID int64 `json:"id"`
}
