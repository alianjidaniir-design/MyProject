package dataModel

import "time"

type User struct {
	ID        int64     `gorm:"column:ID" json:"ID"`
	Code      string    `gorm:"column:code" json:"code"`
	Name      string    `gorm:"column:name" json:"name"`
	Family    string    `gorm:"column:family" json:"family"`
	CreatedAt time.Time `gorm:"column:creatAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
