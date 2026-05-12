package repositories

import (
	"MyProject/apiSchema/InventorySchema"
	"MyProject/apiSchema/commonSchema"
	"context"
)

type InventoryRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[InventorySchema.CreateInventory]) (res InventorySchema.DetailInventory, errStr string, code int, err error)
}

var InventoryRepo InventoryRepository
