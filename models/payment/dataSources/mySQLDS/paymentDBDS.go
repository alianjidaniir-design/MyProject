package mySQLDS

import (
	"MyProject/apiSchema/paymentSchema"
	"MyProject/models/payment/dataModels"
	"MyProject/models/payment/dataSources"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PaymentDBDS struct {
	tableName string
	db        *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func NewPaymentDBDS(tableName string, db *sql.DB) (dataSources.PaymentDS, error) {

	paymentDBInstance := &PaymentDBDS{
		tableName: tableName,
		db:        db,
	}
	return paymentDBInstance, nil
}

func (ds *PaymentDBDS) CreatePayment(ctx context.Context, req paymentSchema.ConfirmationSchema) (dataModels.Payment, error) {
	var payment dataModels.Payment
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return dataModels.Payment{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	var checkTuition bool
	checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM tuition WHERE row = ? AND deleted_at IS NULL) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, checkQuery, req.TuitionRow).Scan(&checkTuition)
	if err != nil {
		return dataModels.Payment{}, err
	}
	if !checkTuition {
		return dataModels.Payment{}, errors.New("tuition does not exist")
	}

	var cash int
	var credited int
	selectQuery := fmt.Sprintf("SELECT debit_amount FROM tuition WHERE row = ?")
	err = tx.QueryRowContext(ctx, selectQuery, req.TuitionRow).Scan(&cash)
	if err != nil {
		return dataModels.Payment{}, errors.New("error in select from tuition")
	}
	if req.Operation == true {
		if req.PaymentType == constants.Cash && req.Operation == true {
			payment.CashAmount = dataModels.NullInt64{Val: int64(cash)}
			payment.NumberInstallment = dataModels.NullString{Val: "0"}
			credited = int(payment.CashAmount.Val)

		} else if req.PaymentType == constants.Installment && req.Operation == true {
			number, err := strconv.Atoi(req.NumberInstallment)
			if err != nil {
				return dataModels.Payment{}, err
			}
			payment.NumberInstallment = dataModels.NullString{Val: req.NumberInstallment}
			payment.InstallmentTotal = dataModels.NullInt64{Val: int64(cash)}
			payment.InstallmentAmount = dataModels.NullInt64{Val: int64(cash / number)}

			credited = int(payment.InstallmentTotal.Val)

		} else {
			return dataModels.Payment{}, errors.New("invalid payment type")
		}
	} else {
		payment.NumberInstallment = dataModels.NullString{Val: "0"}
		payment.InstallmentTotal = dataModels.NullInt64{Val: 0}
		payment.InstallmentAmount = dataModels.NullInt64{Val: 0}
		payment.CashAmount = dataModels.NullInt64{Val: 0}
		credited = 0
	}

	now := time.Now().In(myLocation())
	insertQuery := fmt.Sprintf("INSERT INTO %s (tuition_row , payment_type , number_installment , installment_total , installment_amount , cash_amount , bank , Operation , created_at , updated_at) VALUES (?, ?, ?, ? , ? , ? ,? , ? , ? , ?)", ds.tableName)

	row, err := tx.ExecContext(ctx, insertQuery, req.TuitionRow, req.PaymentType, payment.NumberInstallment.Val, payment.InstallmentTotal.Val, payment.InstallmentAmount.Val, payment.CashAmount.Val, req.Bank, req.Operation, now, now)
	if err != nil {
		return dataModels.Payment{}, err
	}
	lastInsertId, err := row.LastInsertId()
	if err != nil {
		return dataModels.Payment{}, err
	}

	query := `
UPDATE tuition
SET credit_amount = credit_amount + ?
WHERE row = ? AND credit_amount = 0 
`
	res, err := tx.ExecContext(ctx, query, credited, req.TuitionRow)
	if err != nil {
		return dataModels.Payment{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return dataModels.Payment{}, err
	}
	if affected == 0 && req.Operation == true {
		return dataModels.Payment{}, errors.New("credit_amount is not zero (or row not found), update skipped")
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Payment{}, errors.New("the transaction is unsuccessful")
	}
	return ds.getPaymentStudent(ctx, lastInsertId)
}

func (ds *PaymentDBDS) DeletePayment(ctx context.Context, req paymentSchema.DeleteInformation) (res dataModels.Payment, err error) {
	err = ds.checkPayment(ctx, req.ID)
	if err != nil {
		return dataModels.Payment{}, err
	}
	pay, err := ds.getPaymentStudent(ctx, req.ID)
	if err != nil {
		return dataModels.Payment{}, err
	}
	if pay.Operation == true {
		return dataModels.Payment{}, errors.New("You cannot delete this successful payment ")
	}
	now := time.Now().In(myLocation())
	deleteQuery := fmt.Sprintf("UPDATE %s SET deleted_at = ? , updated_at = ? WHERE ID = ?", ds.tableName)
	row, err := ds.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return dataModels.Payment{}, err
	}
	defer row.Close()
	_, err = row.ExecContext(ctx, now, now, req.ID)
	if err != nil {
		return dataModels.Payment{}, err
	}
	return ds.getPaymentStudent(ctx, req.ID)

}

func (ds *PaymentDBDS) DetailPayment(ctx context.Context, req paymentSchema.GetInformation) (res dataModels.Payment, err error) {
	err = ds.checkPayment(ctx, req.ID)
	if err != nil {
		return dataModels.Payment{}, err
	}
	return ds.getPaymentStudent(ctx, req.ID)
}

func (ds *PaymentDBDS) ListPayment(ctx context.Context, req paymentSchema.ListPayment) (res []dataModels.Payment, err error) {
	var payments []dataModels.Payment
	var rows *sql.Rows
	if req.Filter == nil {
		selectQuery := fmt.Sprintf("SELECT * FROM %s ", ds.tableName)
		rows, err = ds.db.QueryContext(ctx, selectQuery)
		if err != nil {
			return nil, err
		}
	} else {
		filterQuery, args, err := ds.filterQuery(*req.Filter)
		if err != nil {
			return nil, err
		}
		rows, err = ds.db.QueryContext(ctx, filterQuery, args...)
		if err != nil {
			return nil, err
		}

	}
	defer rows.Close()
	for rows.Next() {
		var payment dataModels.Payment
		var createdAt, updatedAt, deletedAt sql.NullTime
		var PaymentTotal, PaymentCash, PaymentAmount dataModels.NullInt64
		var PaymentNumber dataModels.NullString
		err = rows.Scan(&payment.ID, &payment.TuitionRow, &payment.PaymentType, &PaymentNumber.Val, &PaymentTotal.Val, &PaymentAmount.Val, &PaymentCash.Val, &payment.Bank, &payment.Operation, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			return nil, err
		}
		payment.NumberInstallment = PaymentNumber
		payment.InstallmentTotal = PaymentTotal
		payment.InstallmentAmount = PaymentAmount
		payment.CashAmount = PaymentCash
		if createdAt.Valid {
			payment.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			payment.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			payment.DeletedAt = deletedAt.Time
		}
		payments = append(payments, payment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (ds *PaymentDBDS) filterQuery(req paymentSchema.Filter) (string, []interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ", ds.tableName)
	var args []interface{}
	condition := []string{}

	if req.PaymentType == "cash" || req.PaymentType == "installment" {
		condition = append(condition, "payment_type = ?")
		args = append(args, req.PaymentType)
	} else {
		return "", nil, errors.New("invalid payment type")
	}
	banks := []string{"meli", "melat", "saderat"}

	check := false
	for _, bank := range banks {

		if req.Bank == bank {
			condition = append(condition, "bank = ?")
			args = append(args, req.Bank)
			check = true
			break
		}
	}
	if !check {
		invalid := "This bank does not exist."
		return invalid, nil, errors.New(invalid)
	}
	if req.Operation != false {
		condition = append(condition, "operation = ?")
		args = append(args, req.Operation)
	} else if req.Operation == true {
		condition = append(condition, "operation = ?")
		args = append(args, req.Operation)
	}
	if len(condition) > 0 {
		query += strings.Join(condition, " AND ")
	}

	query += " ORDER BY id "
	return query, args, nil
}

func (ds *PaymentDBDS) getPaymentStudent(ctx context.Context, ID int64) (dataModels.Payment, error) {
	var payment dataModels.Payment
	var createdAt, updatedAt, deletedAt sql.NullTime
	var PaymentTotal, PaymentCash, PaymentAmount dataModels.NullInt64
	var PaymentNumber dataModels.NullString
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE ID = ?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&payment.ID, &payment.TuitionRow, &payment.PaymentType, &PaymentNumber.Val, &PaymentTotal.Val, &PaymentAmount.Val, &PaymentCash.Val, &payment.Bank, &payment.Operation, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return dataModels.Payment{}, err
	}
	payment.NumberInstallment = PaymentNumber
	payment.InstallmentTotal = PaymentTotal
	payment.InstallmentAmount = PaymentAmount
	payment.CashAmount = PaymentCash
	if createdAt.Valid {
		payment.CreatedAt = createdAt.Time
	}

	if updatedAt.Valid {
		payment.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		payment.DeletedAt = deletedAt.Time
	}
	return payment, nil
}

func (ds *PaymentDBDS) checkPayment(ctx context.Context, ID int64) error {
	var check bool
	searchQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM payments WHERE ID = ? AND deleted_at IS NULL) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, searchQuery, ID).Scan(&check)
	if err != nil {
		return errors.New(err.Error())
	}
	if !check {
		return errors.New("This is not a valid ID")
	}
	return nil

}
