package dataModels

import (
	"database/sql/driver"
	"encoding/json"
)

func (n NullInt64) Value() (driver.Value, error) {
	if n.Val == 0 {
		return nil, nil
	}

	return n.Val, nil
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Val == 0 {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}

func (n NullString) Value() (driver.Value, error) {
	if n.Val == "" {
		return nil, nil
	}

	return n.Val, nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Val == "" {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}
