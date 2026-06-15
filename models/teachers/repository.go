package teachers

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/teacherSchema"
	"MyProject/midddleware/authz"
	"MyProject/models/teachers/dataSources"
	mysqlDataSource "MyProject/models/teachers/dataSources/mysqlDS"
	"MyProject/models/teachers/dataSources/redis"
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
	DBDS     dataSources.TeacherDS
	RedisDS  dataSources.RedisDS
	initRepo error
}

var (
	once    sync.Once
	repoIns *Repository
)

func initIns() {
	DBConn, err := mysqlDataSource.LoadConfig()
	if err != nil {
		return
	}
	open, err := mysqlDataSource.Open(DBConn)
	if err != nil {
		return
	}
	newTeacher, err := mysqlDataSource.NewTeacherDBDS(DBConn.TeacherTableName, open)
	if err != nil {
		repoIns = &Repository{initRepo: fmt.Errorf("error in newTeacher")}
		return
	}
	red, err := redis.NewRedisDS(configs.Addr, configs.Password)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Printf("Warning: Continuing without Redis - blacklist feature will be disabled")
		repoIns = &Repository{initRepo: fmt.Errorf("error in newTeacher")}
	}
	repoIns = &Repository{DBDS: newTeacher, RedisDS: red, initRepo: nil}
	log.Println("success in newTeacher")
}

func GetRepo() *Repository {
	once.Do(initIns)
	return repoIns
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.InformationSchema]) (res teacherSchema.TeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.TeacherSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.TeacherSchema{}, "02", status.StatusUnauthorized, err
	}
	create, err := repo.db().CreateTeacher(ctx, req.Body)

	if err != nil {
		return teacherSchema.TeacherSchema{}, "02", status.UnAvailableServiceError, err
	}
	return teacherSchema.TeacherSchema{Teacher: create}, "04", status.StatusOK, err
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.PaginationSchema]) (res teacherSchema.ListSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.ListSchema{}, "01", 0, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.ListSchema{}, "02", status.StatusUnauthorized, errors.New("wrong db connection")
	}
	list, total, err := repo.db().ListTeachers(ctx, req.Body)
	if err != nil {
		return teacherSchema.ListSchema{}, "03", status.UnAvailableServiceError, err
	}
	return teacherSchema.ListSchema{Teachers: list, Total: total}, "04", 0, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.GetTeacherSchema]) (res teacherSchema.DetailTeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.DetailTeacherSchema{}, "01", 0, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.DetailTeacherSchema{}, "02", 0, errors.New("wrong db connection")
	}
	get, err := repo.db().GetTeacherById(ctx, req.Body)
	if err != nil {
		return teacherSchema.DetailTeacherSchema{}, "03", status.StatusUnauthorized, err
	}
	return teacherSchema.DetailTeacherSchema{Teacher: get}, "", status.StatusOK, nil

}

func (repo *Repository) SoftDelete(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]) (res teacherSchema.SoftDeleteTeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.SoftDeleteTeacherSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.SoftDeleteTeacherSchema{}, "02", status.StatusBadRequest, err
	}
	softDelete, err := repo.db().SoftDeleteTeachers(ctx, req.Body)
	if err != nil {
		return teacherSchema.SoftDeleteTeacherSchema{}, "03", status.StatusUnauthorized, err
	}
	return teacherSchema.SoftDeleteTeacherSchema{Teacher: softDelete, Massage: "teacher deleted successfully"}, "04", 0, nil
}
func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]) (res teacherSchema.UpdateTeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.UpdateTeacherSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.UpdateTeacherSchema{}, "02", status.StatusBadRequest, errors.New("wrong db connection")
	}
	updated, err := repo.db().UpdateTeachers(ctx, req.Body)
	if err != nil {
		return teacherSchema.UpdateTeacherSchema{}, "03", status.StatusUnauthorized, err
	}
	return teacherSchema.UpdateTeacherSchema{Teacher: updated}, "", status.StatusOK, nil
}

func (repo *Repository) Login(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.LoginTeacherRequest], c *fiber.Ctx) (res teacherSchema.EntryStudentSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.EntryStudentSchema{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.db() == nil {
		return teacherSchema.EntryStudentSchema{}, "02", status.StatusInternalServerError, errors.New("wrong db connection")
	}
	access, ref, massage, err := repo.db().LoginTeachers(ctx, req.Body)
	if err != nil {
		return teacherSchema.EntryStudentSchema{}, "03", status.StatusUnauthorized, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    ref,
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

	return teacherSchema.EntryStudentSchema{Massage: massage}, "", status.StatusOK, nil
}

func (repo *Repository) RefreshToken(ctx context.Context, c *fiber.Ctx) (errStr string, code int, err error) {
	if repo.initRepo != nil {
		return "", status.StatusUnauthorized, repo.initRepo
	}
	if repo.db() == nil {
		return "", status.StatusInternalServerError, errors.New("wrong db connection")
	}
	oldRefreshToken, err := dataSource.Cookies(c)
	if err != nil {
		return "", status.StatusUnauthorized, err
	}
	accessToken, refreshToken, err := repo.db().Refresh(ctx, oldRefreshToken)
	if err != nil {
		fmt.Println(oldRefreshToken)
		return "", status.StatusUnauthorized, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
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

func (repo *Repository) Logout(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.LogoutSchema], c *fiber.Ctx) (res teacherSchema.EntryStudentSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.EntryStudentSchema{}, "", status.StatusUnauthorized, repo.initRepo
	}
	if repo.db() == nil {
		return teacherSchema.EntryStudentSchema{}, "02", status.StatusInternalServerError, errors.New("wrong db connection")
	}
	oldRefreshToken, err := dataSource.Cookies(c)
	if err != nil {
		return teacherSchema.EntryStudentSchema{}, "03", status.StatusUnauthorized, err
	}
	err = repo.db().Logout(ctx, req.Body, oldRefreshToken)
	if err != nil {
		return teacherSchema.EntryStudentSchema{}, "04", status.StatusUnauthorized, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(-time.Hour),
	})
	jtiStr, expTime, err := dataSource.BList(c)
	if err != nil {
		return teacherSchema.EntryStudentSchema{}, "05", status.StatusUnauthorized, err
	}
	err = repo.cache().Logout(ctx, jtiStr, expTime)
	if err != nil {
		return teacherSchema.EntryStudentSchema{}, "06", status.StatusUnauthorized, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().In(timeLoc.MyLocation()).Add(-time.Hour),
	})
	return teacherSchema.EntryStudentSchema{Massage: "logout"}, "", status.StatusOK, nil

}

func (repo *Repository) InfoTeacher(ctx context.Context, c *fiber.Ctx) (res teacherSchema.InfoTeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.InfoTeacherSchema{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.db() == nil {
		return teacherSchema.InfoTeacherSchema{}, "02", status.StatusInternalServerError, errors.New("wrong db connection")
	}
	teacherID := authz.GetUserID(c)
	detail, err := repo.db().MyInfo(ctx, teacherID)
	if err != nil {
		return teacherSchema.InfoTeacherSchema{}, "03", status.StatusUnauthorized, err
	}
	return teacherSchema.InfoTeacherSchema{Teacher: detail}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.TeacherDS {
	return repo.DBDS
}

func (repo *Repository) cache() dataSources.RedisDS {
	return repo.RedisDS
}
