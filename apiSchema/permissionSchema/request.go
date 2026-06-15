package permissionSchema

import "time"

type CreatePermissionReq struct {
	Name          string     `json:"name"`
	TimeValidFrom *time.Time `json:"timeValidFrom"`
	TimeValidUntil  *time.Time `json:"timeValidUntil"`
}

type GetPermissionReq struct {
	ID int64 `json:"id"`
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
