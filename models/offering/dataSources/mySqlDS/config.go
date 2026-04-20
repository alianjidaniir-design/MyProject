package mySqlDS

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultTableName      = "offerings"
	MaxOpenConnections    = 5
	MaxIdleConnections    = 10
	MaxConnectionLifetime = 300
)

type Config struct {
	DSN                   string
	TableName             string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime int
}

func LoadConfig() (Config, error) {
	cfg := Config{
		DSN:                   normalize(strings.TrimSpace(os.Getenv("MYSQL_DSN"))),
		TableName:             strings.TrimSpace(os.Getenv("MYSQL_OFFERING_NAME")),
		MaxOpenConnections:    readEnv(os.Getenv("MYSQL_MAX_OPEN_CONNECTIONS"), MaxOpenConnections),
		MaxIdleConnections:    readEnv(os.Getenv("MYSQL_MAX_IDLE_CONNECTIONS"), MaxIdleConnections),
		MaxConnectionLifetime: readEnv(os.Getenv("MYSQL_MAX_CONNECTION_LIFETIME"), MaxConnectionLifetime),
	}
	if cfg.DSN == "" {
		cfg.DSN = defaultTableName
	}
	if err := ValidateTableName(cfg.TableName); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func normalize(dsn string) string {
	if dsn == "" {
		return ""
	}
	if !strings.Contains(dsn, "?") {
		return dsn + "?parseTime=true&loc=Asia%2FTehran&time_zone=%27%2B03:30%27&charset=utf8mb4"

	}
	base, parseQuery, _ := strings.Cut(dsn, "?")
	queryValues, err := url.ParseQuery(parseQuery)
	if err != nil {
		return ""
	}
	if queryValues.Get("parseTime") == "" {
		queryValues.Set("parseTime", "true")
	}
	if queryValues.Get("loc") == "" {
		queryValues.Set("loc", "Asia/Tehran")
	}
	if queryValues.Get("time_zone") == "" {
		queryValues.Set("time_zone", "'+03:30'")
	}
	if queryValues.Get("charset") == "" {
		queryValues.Set("charset", "utf8mb4")
	}
	return fmt.Sprintf("%s%s", base, queryValues.Encode())
}
func readEnv(envKey string, defaultValue int) int {
	raw := strings.TrimSpace(os.Getenv(envKey))
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
