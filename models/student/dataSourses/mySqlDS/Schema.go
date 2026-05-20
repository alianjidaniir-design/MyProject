package mySqlDS

import (
	"fmt"
	"regexp"
)

var safeTableNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func ValidateTableName(tableName string) error {
	if !safeTableNamePattern.MatchString(tableName) {
		return fmt.Errorf("'%s' is invalid %s", tableName, tableName)
	}
	return nil
}
