# راهنمای RBAC در Go (مطابق ساختار پروژه MyProject)

این فایل یک راهنمای یکپارچه برای پیاده‌سازی **Role-Based Access Control (RBAC)** در پروژه فعلی است؛ شامل:

- تعریف Role و Permission
- نگاشت Role -> Permissions
- بررسی مجوز
- Middleware در Fiber
- اتصال به Route های فعلی پروژه

---

## 1) مفهوم Role و Permission

- **Role (نقش):** هویت کاری کاربر (مثل `user`, `editor`, `admin`)
- **Permission (مجوز):** عمل مشخصی که اجازه انجام آن داده می‌شود (مثل `user.delete`, `user.update`)

در RBAC، کاربر مستقیم به عملیات وصل نمی‌شود؛ ابتدا نقش می‌گیرد، و نقش مجموعه‌ای از مجوزها دارد.

---

## 2) ساختار پیشنهادی در همین پروژه

```text
MyProject/
├─ apiSchema/
│  └─ commonSchema/
│     └─ common.go
│
├─ controllers/
│  ├─ mainController/
│  │  └─ main.go
│  └─ user/
│     ├─ create.go
│     ├─ list.go
│     ├─ get.go
│     ├─ update.go
│     ├─ delete.go
│     └─ softDelete.go
│
├─ services/
│  └─ core/
│     ├─ route/
│     │  ├─ route.go
│     │  └─ userRoute.go
│     └─ authz/                  # جدید
│        ├─ rbac.go
│        ├─ middleware.go
│        └─ role_parser.go
│
└─ statics/
   └─ constants/
      ├─ status/
      │  └─ statusCode.go
      ├─ roles/                  # جدید
      │  └─ roles.go
      └─ permissions/            # جدید
         └─ permissions.go
```

---

## 3) چرا این جایگذاری با ساختار فعلی سازگار است؟

- در `apiSchema/commonSchema/common.go` از قبل `Headers` وجود دارد.
- در `controllers/mainController/main.go` با `ParseBody` هدرها پر می‌شوند.
- در `services/core/route` تمام endpoint ها ثبت می‌شوند.
- بنابراین بهترین نقطه اعمال دسترسی: **middleware روی route** است.

نتیجه:
- controller فقط منطق عملیاتی را انجام می‌دهد.
- کنترل دسترسی متمرکز و تمیز می‌ماند.

---

## 4) کدها (Go) - نمونه کامل و قابل استفاده

### 4.1) تعریف نقش‌ها

فایل: `statics/constants/roles/roles.go`

```go
package roles

type Role string

const (
	RoleUser   Role = "user"
	RoleEditor Role = "editor"
	RoleAdmin  Role = "admin"
)
```

### 4.2) تعریف مجوزها

فایل: `statics/constants/permissions/permissions.go`

```go
package permissions

type Permission string

const (
	UserCreate Permission = "user.create"
	UserList   Permission = "user.list"
	UserGet    Permission = "user.get"
	UserUpdate Permission = "user.update"
	UserDelete Permission = "user.delete"
)
```

### 4.3) نگاشت نقش به مجوزها (RBAC Core)

فایل: `services/core/authz/rbac.go`

```go
package authz

import (
	"MyProject/statics/constants/permissions"
	"MyProject/statics/constants/roles"
)

var rolePermissions = map[roles.Role]map[permissions.Permission]struct{}{
	roles.RoleUser: {
		permissions.UserGet: {},
	},
	roles.RoleEditor: {
		permissions.UserGet:    {},
		permissions.UserList:   {},
		permissions.UserUpdate: {},
	},
	roles.RoleAdmin: {
		permissions.UserCreate: {},
		permissions.UserList:   {},
		permissions.UserGet:    {},
		permissions.UserUpdate: {},
		permissions.UserDelete: {},
	},
}

func HasPermission(role roles.Role, permission permissions.Permission) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	_, exists := perms[permission]
	return exists
}
```

### 4.4) تبدیل رشته Role از Header به Role تایپ‌شده

فایل: `services/core/authz/role_parser.go`

```go
package authz

import "MyProject/statics/constants/roles"

func ParseRole(s string) roles.Role {
	switch s {
	case string(roles.RoleAdmin):
		return roles.RoleAdmin
	case string(roles.RoleEditor):
		return roles.RoleEditor
	default:
		return roles.RoleUser
	}
}
```

### 4.5) middleware بررسی مجوز

فایل: `services/core/authz/middleware.go`

```go
package authz

import (
	"MyProject/statics/constants/permissions"
	"MyProject/statics/constants/status"
	"github.com/gofiber/fiber/v2"
)

// فرض فعلی: role در Header با کلید X-Role ارسال می‌شود.
// در نسخه production بهتر است role از JWT claims استخراج شود.
func RequirePermission(p permissions.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleHeader := c.Get("X-Role")
		if roleHeader == "" {
			return c.Status(status.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing role",
				"code":    "AUTH_01",
			})
		}

		role := ParseRole(roleHeader)
		if !HasPermission(role, p) {
			return c.Status(status.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden",
				"code":    "AUTH_02",
			})
		}

		return c.Next()
	}
}
```

### 4.6) استفاده در route های فعلی user

فایل: `services/core/route/userRoute.go`

```go
package route

import (
	. "MyProject/controllers/user"
	"MyProject/services/core/authz"
	"MyProject/statics/constants/permissions"
	"github.com/gofiber/fiber/v2"
)

var userRoute = map[string]string{
	"userCreate":  "/user/create",
	"userList":    "/user/list",
	"userGet":     "/user/get",
	"userUpdate":  "/user/update",
	"userDelete":  "/user/delete",
	"userDelete2": "/user/delete2",
}

func SetupUserRoute(app *fiber.App) map[string]string {
	app.Post(userRoute["userCreate"], authz.RequirePermission(permissions.UserCreate), Create)
	app.Post(userRoute["userList"], authz.RequirePermission(permissions.UserList), List)
	app.Post(userRoute["userGet"], authz.RequirePermission(permissions.UserGet), Get)
	app.Post(userRoute["userUpdate"], authz.RequirePermission(permissions.UserUpdate), Update)
	app.Post(userRoute["userDelete"], authz.RequirePermission(permissions.UserDelete), Delete)
	app.Post(userRoute["userDelete2"], authz.RequirePermission(permissions.UserDelete), SoftDelete)
	return userRoute
}
```

---

## 5) مثال رفتاری

- درخواست `POST /user/delete` با هدر `X-Role: user`:
  - چون `user.delete` ندارد -> پاسخ `403 Forbidden`

- درخواست `POST /user/delete` با هدر `X-Role: admin`:
  - چک مجوز پاس می‌شود -> controller اجرا می‌شود -> repository فراخوانی می‌شود.

---

## 6) نکات عملی برای نسخه حرفه‌ای‌تر

- به‌جای هدر دستی، role را از `JWT Claims` بگیر.
- برای endpoint های بیشتر (course/department/...) مجوزهای جدید تعریف کن.
- اگر نیاز شد map مجوزها را از DB یا config بخوان.
- برای محیط production، کدهای خطای استاندارد پروژه را با `mainController.Error` یکپارچه نگه دار.

---

## 7) جمع‌بندی

در پروژه فعلی، بهترین پیاده‌سازی RBAC:

1. ثابت‌ها در `statics/constants`
2. منطق RBAC در `services/core/authz`
3. اعمال مجوز در middleware روی route
4. controller فقط منطق تجاری را اجرا کند

این الگو هم با ساختار فعلی سازگار است، هم قابل گسترش برای ماژول‌های بعدی.
