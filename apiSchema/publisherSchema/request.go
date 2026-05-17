package publisherSchema

type CreatePublisher struct {
	Name    string `json:"name" validate:"required,max=63,alpha"`
	Phone   string `json:"phone" validate:"required,len=11,numeric"`
	Address string `json:"address" validate:"required,max=255"`
}

type GetPublisher struct {
	ID int64 `json:"id" validate:"required"`
}

type PaginationPublisher struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"required"`
}
