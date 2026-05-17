package publisherSchema

type CreatePublisher struct {
	Name    string `json:"name" val:"required,max=63,alpha"`
	Phone   int    `json:"phone" val:"required,len=11,numeric"`
	Address string `json:"address" val:"required,max=255"`
}
