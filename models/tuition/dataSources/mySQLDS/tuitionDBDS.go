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
	var check bool
	studentQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE student_id = ?) THEN 1 ELSE 0 END
`
	err = tx.QueryRow(studentQuery, req.StudentID).Scan(&check)
	if err != nil {
		return dataModels.Tuition{}, errors.New("student not exist or not enrolled")
	}
	var dbOffering any

	if req.OfferingRow != 0 {
		var checking bool
		dbOffering = req.OfferingRow
		checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE offering_row = ? AND status = 'enrolled' AND deleted_at IS NULL AND student_id = ?) THEN 1 ELSE 0 END
`

		err = tx.QueryRow(checkQuery, dbOffering, req.StudentID).Scan(&checking)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		if !checking {
			return dataModels.Tuition{}, errors.New("offering exist or not enrolled")
		}
	} else {
		dbOffering = nil
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

	req.FixedTuition = constants.FixedTuition

	if req.CourseTuition != 0 && dbOffering != nil {
		req.FixedTuition = 0
		totalDebit = req.CourseTuition
		if req.ExtraOption != 0 {
			totalDebit += req.ExtraOption
		}
	} else if req.CourseTuition == 0 && dbOffering == nil {
		req.CourseTuition = 0
		totalDebit = constants.FixedTuition
		fix := constants.FixedTuition
		var number int
		counted := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE fixed_tuition = ? AND student_id = ? ", ds.tableName)
		err = tx.QueryRowContext(ctx, counted, fix, req.StudentID).Scan(&number)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		if number >= 1 {
			return dataModels.Tuition{}, errors.New(" fixed tuition exists already")
		}
		req.FixedTuition = fix
	} else {
		return dataModels.Tuition{}, errors.New("this kind is invalid")
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

func (ds *TuitionDBDS) UpdateTuition(ctx context.Context, req tuitionSchema.UpdateTuition) (res dataModels.Tuition, err error) {
	err = ds.checkTuition(ctx, req.Row)
	if err != nil {
		return dataModels.Tuition{}, err
	}
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

	now := time.Now().In(myLocation())
	currentTuition, err := ds.selectTuitionByID(ctx, req.Row)
	if err != nil {
		return dataModels.Tuition{}, err
	}

	if currentTuition.FixedTuition == 0 {
		debit := req.CourseTuition + req.ExtraOption
		updated := fmt.Sprintf("UPDATE %s SET  course_tuition = ? , extra_option = ? , debit_amount = ? , updated_at = ?  WHERE row = ?", ds.tableName)
		rows, err := tx.PrepareContext(ctx, updated)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		defer rows.Close()
		_, err = rows.ExecContext(ctx, req.CourseTuition, req.ExtraOption, debit, now, req.Row)
		if err != nil {
			return dataModels.Tuition{}, err
		}

	} else {
		return dataModels.Tuition{}, errors.New("course and option tuition is zero")
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("Error updating tuition: %s", err)
	}
	return ds.selectTuitionByID(ctx, req.Row)
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

func (ds *TuitionDBDS) checkTuition(ctx context.Context, ID int64) error {
	var ok bool
	selectQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM tuition WHERE row = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&ok)
	if err != nil {
		return fmt.Errorf("Error checking tuition existence: %w", err)
	}
	if !ok {
		return errors.New("tuition not exist")
	}
	return nil
}
