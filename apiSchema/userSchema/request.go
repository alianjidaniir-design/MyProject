package userSchema

type LoginRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Family string `json:"family"`
}

type ListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetRequest struct {
	ID int64 `json:"ID"`
}

type UpdateUserRequest struct {
	ID int64 `json:"ID"`
}
