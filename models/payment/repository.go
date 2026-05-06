package payment

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/paymentSchema"
	"MyProject/models/payment/dataSources"
	"MyProject/models/payment/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.PaymentDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepo() {
	cfg, err := mySQLDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: errors.New("Loading config failed . please checking again")}
		return
	}
	open, err := mySQLDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: errors.New("Opening config failed . please checking again")}
	}
	newPayment, err := mySQLDS.NewPaymentDBDS(cfg.PaymentTableName, open)
	if err != nil {
		repo = &Repository{
			initRepo: errors.New("Creating new payment failed . please checking again"),
		}
	}
	repo = &Repository{DBDS: newPayment}
	log.Printf("successfully init repo")

}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.ConfirmationSchema]) (res paymentSchema.DetailPaymentSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return paymentSchema.DetailPaymentSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return paymentSchema.DetailPaymentSchema{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreatePayment(ctx, req.Body)
	if err != nil {
		return paymentSchema.DetailPaymentSchema{}, "03", status.StatusBadRequest, err
	}
	return paymentSchema.DetailPaymentSchema{Detail: create}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.DeleteInformation]) (res paymentSchema.DetailChangePaymentSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return paymentSchema.DetailChangePaymentSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return paymentSchema.DetailChangePaymentSchema{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeletePayment(ctx, req.Body)
	if err != nil {
		return paymentSchema.DetailChangePaymentSchema{}, "03", status.StatusBadRequest, err
	}
	return paymentSchema.DetailChangePaymentSchema{Detail: deleted, Massage: "payment deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.GetInformation]) (res paymentSchema.DetailPaymentSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return paymentSchema.DetailPaymentSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return paymentSchema.DetailPaymentSchema{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().DetailPayment(ctx, req.Body)
	if err != nil {
		return paymentSchema.DetailPaymentSchema{}, "03", status.StatusBadRequest, err
	}
	return paymentSchema.DetailPaymentSchema{Detail: get}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[paymentSchema.ListPayment]) (res paymentSchema.DetailListPaymentSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return paymentSchema.DetailListPaymentSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return paymentSchema.DetailListPaymentSchema{}, "02", status.StatusInternalServerError, err
	}
	list, err := repo.db().ListPayment(ctx, req.Body)
	if err != nil {
		return paymentSchema.DetailListPaymentSchema{}, "03", status.StatusBadRequest, err
	}
	return paymentSchema.DetailListPaymentSchema{Detail: list}, "", status.StatusOK, nil

}

func (repo *Repository) db() dataSources.PaymentDS {
	return repo.DBDS
}
