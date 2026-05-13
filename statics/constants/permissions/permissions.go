package permissions

type Permissions string

const (
	CreateUser Permissions = "userCreate"
)

var Permission = struct {
	UserCreate string
	UserRead   string
}{
	UserCreate: "UserCreate",
	UserRead:   "UserRead",
}
