package mysqlDS

import (
	"database/sql"
	"time"
)

func Open(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxConnectionLifetime) * time.Second)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
