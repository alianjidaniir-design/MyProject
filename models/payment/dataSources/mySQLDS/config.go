package mySqlDS

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultTableName      = "courses"
	MaxOpenConnections    = 10
	MaxIdleConnections    = 5
	MaxConnectionLifetime = 300
)

type Config struct {
	DSN                   string
	CourseTableName       string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime int
}

func LoadConfig() (Config, error) {
	cfg := Config{
		DSN:                   normalize(strings.TrimSpace(os.Getenv("MYSQL_DSN"))),
		CourseTableName:       strings.TrimSpace(os.Getenv("MYSQL_COURSE_TABLE_NAME")),
		MaxOpenConnections:    readEnvInt("MYSQL_MAX_OPEN_CONNECTIONS", MaxOpenConnections),
		MaxIdleConnections:    readEnvInt("MYSQL_MAX_IDLE_CONNECTIONS", MaxIdleConnections),
		MaxConnectionLifetime: readEnvInt("MYSQL_MAX_CONNECTION_LIFETIME", MaxConnectionLifetime),
	}
	if cfg.CourseTableName == "" {
		cfg.CourseTableName = defaultTableName
	}

	if err := validateTableName(cfg.CourseTableName); err != nil {
		return Config{}, err
	}
	return cfg, nil

}

func normalize(dsn string) string {
	if dsn == "" {
		return ""
	}
	if strings.Contains(dsn, "?") {
		return dsn + "?parseTime=true&loc=Asia%2FTehran&time_zone=%27%2B03:30%27&charset=utf8mb4"
	}
	base, parseQuery, _ := strings.Cut(dsn, "?")
	queryValues, err := url.ParseQuery(parseQuery)
	if err != nil {
		return dsn
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
	return fmt.Sprintf("%s?%s", base, queryValues.Encode())

}

func readEnvInt(envKey string, defaultValue int) int {
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
