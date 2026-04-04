package dataModel

type User struct {
	ID     int64  `gorm:"column:ID" json:"ID"`
	Code   string `gorm:"column:code" json:"code"`
	Name   string `gorm:"column:name" json:"name"`
	Family string `gorm:"column:family" json:"family"`
}
