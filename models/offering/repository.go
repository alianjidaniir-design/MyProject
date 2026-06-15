package offering

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/offeringSchema"
	"MyProject/midddleware/authz"
	"MyProject/models/offering/dataSources"
	"MyProject/models/offering/dataSources/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	DBDS     dataSources.OfferingDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepository() {
	cfg, err := mySqlDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: errors.New("you can not load config" + err.Error())}
		return
	}
	open, err := mySqlDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: errors.New("you can not open config" + err.Error())}
		return
	}
	newOffer, err := mySqlDS.NewOfferingDBDS(cfg.TableName, open)
	if err != nil {
		repo = &Repository{initRepo: errors.New(err.Error())}
		return
	}
	repo = &Repository{DBDS: newOffer}
}

func GetRepository() *Repository {
	once.Do(initRepository)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.CreateOfferingRequest]) (res offeringSchema.CreateOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.CreateOfferingResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.CreateOfferingResponse{}, "02", status.StatusBadRequest, err
	}
	create, err := repo.db().CreateOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.CreateOfferingResponse{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.CreateOfferingResponse{Specification: create}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.ListOfferingsRequest]) (res offeringSchema.ListOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.ListOfferingResponse{}, "04", status.StatusBadRequest, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.ListOfferingResponse{}, "05", status.StatusBadRequest, err
	}
	list, total, err := repo.db().ListOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.ListOfferingResponse{}, "06", status.StatusInternalServerError, err
	}
	return offeringSchema.ListOfferingResponse{Offerings: list, TotalCount: total}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.GetRowOfferingRequest]) (res offeringSchema.DetailOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.DetailOfferingResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.DetailOfferingResponse{}, "02", status.StatusBadRequest, err
	}
	get, err := repo.db().GetOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.DetailOfferingResponse{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.DetailOfferingResponse{Specification: get}, "", status.StatusOK, nil
}

func (repo *Repository) DeActive(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.GetRowOfferingRequest]) (res offeringSchema.DeactivateOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.DeactivateOfferingResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.DeactivateOfferingResponse{}, "02", status.StatusBadRequest, err
	}
	deActive, err := repo.db().DeActiveOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.DeactivateOfferingResponse{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.DeactivateOfferingResponse{Specification: deActive}, "", status.StatusOK, nil

}
func (repo *Repository) Edit(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.EditOffering]) (res offeringSchema.ViewAfterEditCourse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.ViewAfterEditCourse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.ViewAfterEditCourse{}, "02", status.StatusBadRequest, err
	}
	edit, err := repo.db().EditOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.ViewAfterEditCourse{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.ViewAfterEditCourse{Massage: "View offering After Edition", Specification: edit}, "", status.StatusOK, nil

}

func (repo *Repository) ListClassesTeacher(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.Pages], c *fiber.Ctx) (res offeringSchema.ClassesTeacher, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.ClassesTeacher{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.ClassesTeacher{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	teacherID := authz.GetUserID(c)
	detail, tot, page, err := repo.db().ListClassesTeacher(ctx, req.Body, teacherID)
	if err != nil {
		return offeringSchema.ClassesTeacher{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.ClassesTeacher{MyClasses: detail, Total: tot, Page: page}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.OfferingDS {
	return repo.DBDS
}
