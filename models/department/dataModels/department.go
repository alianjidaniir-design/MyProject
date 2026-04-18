package dataModels

import "time"

type Department struct {
	ID               int64     `json:"id"`
	College          string    `json:"college"`
	EducationalGroup string    `json:"educational_group"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
