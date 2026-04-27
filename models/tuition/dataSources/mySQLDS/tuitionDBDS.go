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
	var dbCourse any

	if req.CourseID != 0 {
		var check bool
		dbCourse = req.CourseID
		checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE course_id = ? AND student_id = ?) THEN 1 ELSE 0 END
`
		err = tx.QueryRow(checkQuery, req.CourseID, req.StudentID).Scan(&check)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		if !check {
			fmt.Println(req.StudentID, req.CourseID)
			return dataModels.Tuition{}, errors.New("student or course not exist")
		}
	} else {
		var check bool
		dbCourse = nil
		checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE student_id = ?) THEN 1 ELSE 0 END
`
		err = tx.QueryRow(checkQuery, req.StudentID).Scan(&check)
		if err != nil {
			return dataModels.Tuition{}, err
		}
		if !check {
			fmt.Println(req.StudentID, req.CourseID)
			return dataModels.Tuition{}, errors.New("student not exist")
		}
	}

	var lastID int64

	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row), 0) FROM %s", ds.tableName)
	err = tx.QueryRowContext(ctx, lastIDQuery).Scan(&lastID)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("failed to get last tuition row: %w", err)
	}

	newID := lastID + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , student_id, course_id , fixed_tuition , course_tuition , extra_option , debit_amount  , reminder , created_At , updated_at) VALUES (?, ? , ? , ? , ? , ? , ? , ? , ? , ?)", ds.tableName)
	now := time.Now().In(myLocation())
	var totalDebit int

	if req.ExtraOption != 0 {
		totalDebit += req.ExtraOption

	} else if req.CourseTuition != 0 && req.CourseID != 0 {
		totalDebit += req.CourseTuition
	} else {
		totalDebit = constants.FixedTuition
		req.FixedTuition = constants.FixedTuition
	}
	debitAmount := totalDebit

	_, err = tx.ExecContext(ctx, insertQuery, newID, req.StudentID, dbCourse, req.FixedTuition, req.CourseTuition, req.ExtraOption, debitAmount, now, now)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("Error inserting tuition: %s", err)
	}

	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("خطا در بروزرسانی reminder برای ردیف %v: %w", newID, err)
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Tuition{}, err
	}

	return ds.selectTuitionByID(ctx, newID)

}

func (ds *TuitionDBDS) selectTuitionByID(ctx context.Context, ID int64) (res dataModels.Tuition, err error) {
	var tuition dataModels.Tuition
	var courseID sql.NullInt64

	readQuery := fmt.Sprintf(`
        SELECT row, student_id,course_id, fixed_tuition, course_tuition, extra_option, 	debit_amount ,credit_amount , reminder, created_at, updated_at , deleted_at
        FROM %s
        WHERE row = ? `, ds.tableName)

	var createdAt, updatedAt, deletedAt sql.NullTime
	err = ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&tuition.Row, &tuition.StudentID, &courseID, &tuition.FixedTuition, &tuition.CourseTuition, &tuition.ExtraOption, &tuition.DebitAmount, &tuition.CreditAmount, &tuition.Reminder, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("failed to read tuition by row: %w", err)
	}

	if courseID.Valid {
		tuition.CourseID = courseID.Int64
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
