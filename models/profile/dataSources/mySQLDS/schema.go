package mySQLDS

import (
	"fmt"
	"regexp"
)

var valid = regexp.MustCompile("^[a-zA-Z0-9]+$")

func ValidateTableName(tableName string) error {
	if !valid.MatchString(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}
	return nil
}
