package admins

import (
	"MyProject/apiSchema/adminSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/admins/dataSources"
	mysqlDataSource "MyProject/models/admins/dataSources/mysqlDS"
	"MyProject/models/teachers/dataSources/redis"
	"MyProject/models/token/dataSource"
	"MyProject/statics/configs"
	"MyProject/statics/constants"
	"MyProject/statics/constants/status"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	DBDS     dataSources.AdminDS
	RedisDS  dataSources.RedisDS
	initRepo error
}

var (
	once    sync.Once
	repoIns *Repository
)

func initIns() {
	cfg, err := mysqlDataSource.LoadConfig()
	if err != nil {
		return
	}
	open, err := mysqlDataSource.Open(cfg)
	if err != nil {
		return
	}
	newAdmin, err := mysqlDataSource.NewAdminDBDS(cfg.AdminTableName, open)
	if err != nil {
		repoIns = &Repository{initRepo: fmt.Errorf("error in newAdminD: %w", err)}
		return
	}
	red, err := redis.NewRedisDS(configs.Addr, configs.Password)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Printf("Warning: Continuing without Redis - blacklist feature will be disabled")
		repoIns = &Repository{initRepo: fmt.Errorf("error in newAdmin: %w", err)}
	}
	repoIns = &Repository{DBDS: newAdmin, RedisDS: red, initRepo: nil}
	log.Println("success in newAdmin")
}

func GetRepo() *Repository {
	once.Do(initIns)
	return repoIns
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[adminSchema.InformationSchema]) (res adminSchema.DetailAdminSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return adminSchema.DetailAdminSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.db() == nil {
		return adminSchema.DetailAdminSchema{}, "02", status.StatusUnauthorized, err
	}
	create, err := repo.db().CreateAdmin(ctx, req.Body)

	if err != nil {
		return adminSchema.DetailAdminSchema{}, "03", status.UnAvailableServiceError, err
	}
	return adminSchema.DetailAdminSchema{Detail: create}, "0", status.StatusOK, err
}

func (repo *Repository) Login(ctx context.Context, req commonSchema.BaseRequest[adminSchema.LoginAdminRequest], c *fiber.Ctx) (res adminSchema.EntryAdminSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return adminSchema.EntryAdminSchema{}, "01", status.StatusUnauthorized, err
	}
	if repo.db() == nil {
		return adminSchema.EntryAdminSchema{}, "02", status.StatusUnauthorized, err
	}
	access, refresh, massage, err := repo.db().Login(ctx, req.Body)
	if err != nil {
		return adminSchema.EntryAdminSchema{}, "03", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    access,
		HTTPOnly: true,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(constants.AccessTokenExpiry),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(constants.RefreshTokenExpiry),
	})
	return adminSchema.EntryAdminSchema{Massage: massage}, "", status.StatusOK, nil

}

func (repo *Repository) Refresh(ctx context.Context, c *fiber.Ctx) (errStr string, code int, err error) {
	if repo.initRepo != nil {
		return "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.db() == nil {
		return "02", status.StatusUnauthorized, err
	}
	refresh, err := dataSource.Cookies(c)
	if err != nil {
		return "03", status.UnAvailableServiceError, err
	}
	newAccess, newRef, err := repo.db().Refresh(ctx, refresh)
	if err != nil {
		return "04", status.UnAvailableServiceError, err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    newRef,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(constants.RefreshTokenExpiry),
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    newAccess,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(constants.AccessTokenExpiry),
	})
	return "", status.StatusOK, nil

}

func (repo *Repository) Logout(ctx context.Context, req commonSchema.BaseRequest[adminSchema.LogoutSchema], c *fiber.Ctx) (res adminSchema.EntryAdminSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return adminSchema.EntryAdminSchema{}, "01", status.StatusUnauthorized, err
	}
	if repo.db() == nil {
		return adminSchema.EntryAdminSchema{}, "02", status.StatusUnauthorized, err
	}
	oldRef, err := dataSource.Cookies(c)
	if err != nil {
		return adminSchema.EntryAdminSchema{}, "03", status.UnAvailableServiceError, err
	}
	err = repo.db().Logout(ctx, req.Body, oldRef)
	if err != nil {
		return adminSchema.EntryAdminSchema{}, "04", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	})
	str, tim, err := dataSource.BList(c)
	if err != nil {
		return adminSchema.EntryAdminSchema{}, "05", status.UnAvailableServiceError, err
	}
	err = repo.cache().Logout(ctx, str, tim)
	if err != nil {
		return adminSchema.EntryAdminSchema{}, "06", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	})
	return adminSchema.EntryAdminSchema{Massage: "logout"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.AdminDS {
	return repo.DBDS
}

func (repo *Repository) cache() dataSources.RedisDS {
	return repo.RedisDS
}
