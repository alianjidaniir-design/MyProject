package permissions

type Permissions string

const (
	UserCreate     Permissions = "user/create"
	UserDelete     Permissions = "user/delete"
	UserList       Permissions = "user/list"
	UserUpdate     Permissions = "user/update"
	UserGet        Permissions = "user/get"
	UserSoftDelete Permissions = "user/soft_delete"
)
