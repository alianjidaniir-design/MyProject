package Inventory

import "MyProject/models/inventory/dataModel"

type DetailInventory struct {
	Detail  dataModel.Inventory `json:"detail"`
	Massage string              `json:"massage"`
}
