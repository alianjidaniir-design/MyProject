package dataModel

import "MyProject/models/payment/dataModels"

type Book struct {
	Code       int64                 `json:"code"`
	Name       string                `json:"name"`
	Writer     string                `json:"writer"`
	Translator dataModels.NullString `json:"translator"`
	Publisher  string                `json:"publisher"`
}
