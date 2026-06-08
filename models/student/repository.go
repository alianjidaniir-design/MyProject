package student

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/studentSchema"
	"MyProject/midddleware/authz"
	userDataSourses "MyProject/models/student/dataSources"
	"MyProject/models/student/dataSources/mySqlDS"
	"MyProject/models/student/dataSources/redis"
	"MyProject/models/token/dataSource"
	"MyProject/pkg/timeLoc"
	"MyProject/statics/configs"
	"MyProject/statics/constants"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	dbDS    userDataSourses.StudentDB
	redisDS userDataSourses.RedisDS
	initErr error
}

var (
	once    sync.Once
	repoIns *Repository
)

func initRepoIns() {

	cfg, err := mySqlDS.LoadConfig()
	if err != nil {
		repoIns = &Repository{initErr: fmt.Errorf("failed to load config: %v", err)}
		return
	}
	DBConn, err := mySqlDS.Open(cfg)
	if err != nil {
		repoIns = &Repository{initErr: fmt.Errorf("failed to open config: %v", err)}
		log.Printf("Error opening DB connection: %v", err)
		return
	}

	userDNInstance, err := mySqlDS.NewStudentDBDS(DBConn, cfg.StudentTableName)
	if err != nil {
		_ = DBConn.Close()
		repoIns = &Repository{initErr: fmt.Errorf("failed to connect to DB: %v", err)}
		log.Printf("Error opening DB connection: %v", err)
		return
	}
	red, err := redis.NewRedisDS(configs.Addr, configs.Password)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Printf("Warning: Continuing without Redis - blacklist feature will be disabled")
		repoIns = &Repository{dbDS: userDNInstance, redisDS: nil, initErr: nil}
	} else {

		repoIns = &Repository{dbDS: userDNInstance, redisDS: red, initErr: nil}
	}
	log.Println("repository init success")

}
func GetRepoIns() *Repository {
	once.Do(initRepoIns)
	return repoIns
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[studentSchema.SignUpStudent]) (res studentSchema.DetailStudent, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.DetailStudent{}, "13", status.StatusUnauthorized, repo.initErr
	}
	if repo.dbDS == nil {
		return studentSchema.DetailStudent{}, "14", status.UnAvailableServiceError, errors.New("student dataSource not configured")
	}

	createdUser, err := repo.db().CreateStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.DetailStudent{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.DetailStudent{Student: createdUser}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[studentSchema.ListRequest]) (res studentSchema.ListStudents, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.ListStudents{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.ListStudents{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	lists, total, err := repo.db().ReadStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.ListStudents{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.ListStudents{Students: lists, Total: total}, "", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[studentSchema.UpdateUserRequest]) (res studentSchema.UpdateResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.UpdateResponse{}, "14", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.UpdateResponse{}, "15", status.StatusInternalServerError, errors.New("bad")
	}
	updatedUser, err := repo.db().UpdateStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.UpdateResponse{}, "16", status.UnAvailableServiceError, err
	}
	return studentSchema.UpdateResponse{User: updatedUser}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[studentSchema.GetRequest]) (res studentSchema.GetResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.GetResponse{}, "01", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.GetResponse{}, "02", status.StatusInternalServerError, errors.New("bad")
	}

	getIng, err := repo.db().GetStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.GetResponse{}, "03", status.UnAvailableServiceError, err
	}
	return studentSchema.GetResponse{User: getIng}, "", status.StatusOK, nil
}

func (repo *Repository) SoftDelete(ctx context.Context, req commonSchema.BaseRequest[studentSchema.SoftDeleteRequest]) (res studentSchema.SoftDeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.SoftDeleteResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.SoftDeleteResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	soft, err := repo.db().SoftDeleteStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.SoftDeleteResponse{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.SoftDeleteResponse{User: soft}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[studentSchema.DeleteRequest]) (res studentSchema.DeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.DeleteResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.DeleteResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	deletedUser, err := repo.db().DeleteStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.DeleteResponse{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.DeleteResponse{User: deletedUser}, "", status.StatusOK, nil
}

func (repo *Repository) Entry(ctx context.Context, req commonSchema.BaseRequest[studentSchema.LoginStudent], c *fiber.Ctx) (res studentSchema.StudentEntry, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.StudentEntry{}, "01", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.StudentEntry{}, "02", status.StatusInternalServerError, errors.New("bad")
	}
	access, refresh, massage, err := repo.db().StudentEntry(ctx, req.Body)
	if err != nil {
		return studentSchema.StudentEntry{}, "03", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(constants.RefreshTokenExpiry),
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    access,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(constants.AccessTokenExpiry),
	})

	return studentSchema.StudentEntry{Massage: massage}, "", status.StatusOK, nil

}

func (repo *Repository) RefreshToken(ctx context.Context, c *fiber.Ctx) (errStr string, code int, err error) {
	if repo.initErr != nil {
		return "01", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return "02", status.StatusInternalServerError, errors.New("bad")
	}
	oldRefresh, err := dataSource.Cookies(c)
	if err != nil {
		return "03", status.UnAvailableServiceError, err
	}
	refresh, accessToken, err := repo.db().RefreshToken(ctx, oldRefresh)
	if err != nil {
		return "04", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(constants.RefreshTokenExpiry),
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(constants.AccessTokenExpiry),
	})
	return "", status.StatusOK, nil
}

func (repo *Repository) Logout(ctx context.Context, req commonSchema.BaseRequest[studentSchema.LogoutRequest], c *fiber.Ctx) (res studentSchema.StudentEntry, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.StudentEntry{}, "01", status.StatusUnauthorized, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.StudentEntry{}, "02", status.StatusInternalServerError, errors.New("bad")
	}
	ref, err := dataSource.Cookies(c)
	if err != nil {
		return studentSchema.StudentEntry{}, "03", status.UnAvailableServiceError, err
	}
	err = repo.db().RevokedRefreshToken(ctx, req.Body, ref)
	if err != nil {
		return studentSchema.StudentEntry{}, "04", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})
	jtiStr, expTime, err := dataSource.BList(c)
	if err != nil {
		return studentSchema.StudentEntry{}, "05", status.UnAvailableServiceError, err
	}
	err = repo.cache().Logout(ctx, jtiStr, expTime)
	if err != nil {
		return studentSchema.StudentEntry{}, "06", status.UnAvailableServiceError, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})
	return studentSchema.StudentEntry{Massage: "logout"}, "", status.StatusOK, nil

}

func (repo *Repository) StudentInformation(ctx context.Context, c *fiber.Ctx) (res studentSchema.InfoStudent, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.InfoStudent{}, "01", status.StatusUnauthorized, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.InfoStudent{}, "02", status.StatusInternalServerError, errors.New("bad")
	}
	studentID := authz.GetUserID(c)
	detail, err := repo.db().StudentInformation(ctx, studentID)
	if err != nil {
		return studentSchema.InfoStudent{}, "03", status.UnAvailableServiceError, err
	}
	return studentSchema.InfoStudent{Info: detail}, "", status.StatusOK, nil

}

func (repo *Repository) db() userDataSourses.StudentDB {
	return repo.dbDS
}

func (repo *Repository) cache() userDataSourses.RedisDS {
	if repo.redisDS == nil {
		log.Println("WARNING: redisDS is nil in cache()")
	}
	return repo.redisDS
}
