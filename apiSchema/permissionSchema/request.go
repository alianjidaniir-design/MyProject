package permissionSchema

type CreatePermissionReq struct {
	Name string `json:"name"`
}

type GetPermissionReq struct {
	ID int64 `json:"id"`
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
