package mySQLDS

import (
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition/dataModels"
	tuitionDataSourses "MyProject/models/tuition/dataSources"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TuitionDBDS struct {
	tableName string
	db        *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/ُTehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func NewTuitionDBDS(tableName string, db *sql.DB) (tuitionDataSourses.TuitionDS, error) {

	tuitionDBInstance := &TuitionDBDS{
		tableName: tableName,
		db:        db,
	}
	return tuitionDBInstance, nil
}

func (ds *TuitionDBDS) CreateTuition(ctx context.Context, req tuitionSchema.CreateTuition) (res dataModels.Tuition, err error) {
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return dataModels.Tuition{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	var dbOffering any

	if req.OfferingRow != 0 {
		var check bool
		dbOffering = req.OfferingRow
		checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE offering_row = ? AND status = 'enrolled' AND deleted_at IS NULL AND student_id = ?) THEN 1 ELSE 0 END
`
		err = tx.QueryRow(checkQuery, req.OfferingRow, req.StudentID).Scan(&check)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		if !check {
			return dataModels.Tuition{}, errors.New("student not exist or not enrolled")
		}
	} else {
		return dataModels.Tuition{}, errors.New("offering Row is null")
	}

	var lastID int64

	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row), 0) FROM %s", ds.tableName)
	err = tx.QueryRowContext(ctx, lastIDQuery).Scan(&lastID)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("failed to get last tuition row: %w", err)
	}

	newID := lastID + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , student_id, offering_row , fixed_tuition , course_tuition , extra_option , debit_amount  , created_At , updated_at) VALUES (?, ? , ? , ? , ? , ?  , ? , ? , ?)", ds.tableName)
	now := time.Now().In(myLocation())
	var totalDebit int

	if (!((req.FixedTuition != 0) || (req.CourseTuition != 0 || req.ExtraOption != 0))) || ((req.FixedTuition != 0) && (req.CourseTuition != 0 || req.ExtraOption != 0)) {

		return dataModels.Tuition{}, errors.New("invalid tuition: either fixed tuition must be non-zero, or course/extra tuition must be provided")
	}

	if req.CourseTuition != 0 && dbOffering != nil {
		totalDebit = req.CourseTuition
		fmt.Println(12)
		if req.ExtraOption != 0 {
			totalDebit += req.ExtraOption
		}
	} else {
		req.CourseTuition = 0
		totalDebit = constants.FixedTuition
		fix := constants.FixedTuition
		counted := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE fixed_tuition = ? AND offering_row = ?", ds.tableName)
		err = tx.QueryRowContext(ctx, counted).Scan(&fix, &dbOffering)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		if fix > 1 {
			return dataModels.Tuition{}, errors.New(" fixed tuition exists already")
		}
		req.FixedTuition = fix
	}
	if totalDebit < 0 {
		return dataModels.Tuition{}, errors.New("calculated total debit cannot be negative")
	}

	_, err = tx.ExecContext(ctx, insertQuery, newID, req.StudentID, dbOffering, req.FixedTuition, req.CourseTuition, req.ExtraOption, totalDebit, now, now)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("Error inserting tuition: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		return dataModels.Tuition{}, err
	}

	return ds.selectTuitionByID(ctx, newID)

}

func (ds *TuitionDBDS) selectTuitionByID(ctx context.Context, ID int64) (res dataModels.Tuition, err error) {
	var tuition dataModels.Tuition
	var offeringID sql.NullInt64

	readQuery := fmt.Sprintf(`
        SELECT row, student_id,offering_row, fixed_tuition, course_tuition, extra_option, 	debit_amount ,credit_amount , created_at, updated_at , deleted_at
        FROM %s
        WHERE row = ? `, ds.tableName)

	var createdAt, updatedAt, deletedAt sql.NullTime
	err = ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&tuition.Row, &tuition.StudentID, &offeringID, &tuition.FixedTuition, &tuition.CourseTuition, &tuition.ExtraOption, &tuition.DebitAmount, &tuition.CreditAmount, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("failed to read tuition by row: %w", err)
	}

	if offeringID.Valid {
		tuition.OfferingID = offeringID.Int64
	} else {
		tuition.CourseTuition = 0
	}
	if createdAt.Valid {
		tuition.CreatedAt = createdAt.Time.In(myLocation())
	} else {
		tuition.CreatedAt = time.Time{}
	}

	if updatedAt.Valid {
		tuition.UpdatedAt = updatedAt.Time.In(myLocation())
	} else {
		tuition.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		tuition.DeletedAt = deletedAt.Time.In(myLocation())
	} else {
		tuition.DeletedAt = time.Time{}
	}

	return tuition, nil

}
