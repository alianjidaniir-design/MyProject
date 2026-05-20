package mySqlDS

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)


func Open(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTimeSeconds))
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
