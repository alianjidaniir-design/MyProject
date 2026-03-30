package mySqlDS

import (
	"database/sql"
	"fmt"
	"regexp"
)

var safeTableNamePattern = regexp.MustCompile(`[^a-zA-Z0-9_]+$`)

func ValidateTableName(tableName string) error {
	if !safeTableNamePattern.MatchString(tableName) {
		return fmt.Errorf("table name '%s' is invalid", tableName)
	}
	return nil
}

func studentTableName(tableName string) (string, error) {
	if err := ValidateTableName(tableName); err != nil {
		return "", nil
	}
	return fmt.Sprintf("`%s`", tableName), nil
}

func EnsureStudentTable(db *sql.DB, tableName string) error {
	tableIdentifier, err := studentTableName(tableName)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    family VARCHAR(120) NOT NULL
);`, tableIdentifier)
	_, err = db.Exec(query)
	return err
}
