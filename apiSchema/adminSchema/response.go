package adminSchema

import (
	adminDataModels "MyProject/models/admins/dataModels"
)

type AdminSchema struct {
	Admin adminDataModels.Admins `json:"admin"`
}

type ListSchema struct {
	Admins []adminDataModels.Admins `json:"admins"`
	Total  int64                    `json:"total"`
}

type DetailAdminSchema struct {
	Detail adminDataModels.Admins `json:"detail"`
}

type HardDeleteAdminSchema struct {
	Massage string `json:"massage"`
}

type SoftDeleteAdminSchema struct {
	Admin   adminDataModels.Admins `json:"admin"`
	Massage string                 `json:"massage"`
}

type UpdateAdminSchema struct {
	Admin adminDataModels.Admins `json:"admin"`
}

type EntryAdminSchema struct {
	Massage string `json:"massage"`
}
