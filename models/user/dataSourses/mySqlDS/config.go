package mySqlDS

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultsStudentTableName      = "student"
	defaultMaxOpenConnections     = 10
	defaultMaxIdleConnections     = 5
	defaultConnMaxLifeTimeSeconds = 300
)

type Config struct {
	DSN                    string
	StudentTableName       string
	MaxOpenConnections     int
	MaxIdleConnections     int
	ConnMaxLifeTimeSeconds int
}

func LoadConfig() (cfg Config, err error) {
	cfg = Config{
		DSN:                    normalizeDSN(strings.TrimSpace(os.Getenv("MYSQL_DSN"))),
		StudentTableName:       strings.TrimSpace(os.Getenv("MYSQL_STUDENT_TABLE")),
		MaxOpenConnections:     readEnvInt("MYSQL_MAX_OPEN_CONNECTIONS", defaultMaxOpenConnections),
		MaxIdleConnections:     readEnvInt("MYSQL_MAX_IDLE_CONNECTIONS", defaultMaxIdleConnections),
		ConnMaxLifeTimeSeconds: readEnvInt("MYSQL_CONN_MAX_LIFE_TIME_SECONDS", defaultConnMaxLifeTimeSeconds),
	}
	if cfg.StudentTableName == "" {
		cfg.StudentTableName = defaultsStudentTableName
	}

	if err := ValidateTableName(cfg.StudentTableName); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func normalizeDSN(dsn string) string {
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
