package publisherSchema

type CreatePublisher struct {
	Name    string `json:"name" validate:"required,max=63,alpha"`
	Phone   int    `json:"phone" validate:"required,len=11,numeric"`
	Address string `json:"address" validate:"required,max=255"`
}
