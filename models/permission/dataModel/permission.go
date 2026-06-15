package dataModel

import "time"

type Permission struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	TimeValidFrom *time.Time `json:"timeValidFrom"`
	TimeValidUntil  *time.Time `json:"timeValidUntil"`
}
