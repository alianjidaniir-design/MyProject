package mySQLDS

import (
	"database/sql"
)

type InventoryDBDS struct {
	tableName string
	db        *sql.DB
}

func NewInventoryDBDS(tableName string, db *sql.DB) error {
	return nil
}
