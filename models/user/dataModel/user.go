package dataModel

type User struct {
	ID     int64  `gorm:"column:ID" json:"ID" msgpack:"ID"`
	Code   string `gorm:"column:code" json:"code" msgpack:"code"`
	Name   string `gorm:"column:name" json:"name" msgpack:"name"`
	Family string `gorm:"column:family" json:"family" msgpack:"family"`
}
