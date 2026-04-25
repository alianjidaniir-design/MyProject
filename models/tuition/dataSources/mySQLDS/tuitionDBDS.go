package mySQLDS

import (
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition/dataModels"
	tuitionDataSourses "MyProject/models/tuition/dataSources"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TuitionDBDS struct {
	tableName string
	tableSQL  string
	db        *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/ُTehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func NewTuitionDBDS(db *sql.DB, tableName string) (tuitionDataSourses.TuitionDS, error) {

	tuitionDBInstance := &TuitionDBDS{
		tableName: tableName,
		tableSQL:  tableName,
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

	var checkStudent, checkCourse bool
	checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM  registration WHERE student_id = ? ) THEN 1 ELSE 0 END,
CASE WHEN EXISTS (SELECT 1 FROM  registration WHERE course_id = ? ) THEN 1 ELSE 0 END
`
	err = tx.QueryRow(checkQuery, req.StudentID, req.CourseID).Scan(&checkStudent, &checkCourse)
	if err != nil {
		return dataModels.Tuition{}, err
	}
	if !checkStudent || !checkCourse {
		return dataModels.Tuition{}, errors.New("student or course not exist")
	}
	var lastID int64

	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row), 0) FROM %s", ds.tableName)
	err = tx.QueryRowContext(ctx, lastIDQuery).Scan(&lastID)
	if err != nil {
		return dataModels.Tuition{}, err
	}
	newID := lastID + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , student_id, course_id , fixed_tuition , course_tuition , extra_option , debit_amount , credit_amount , reminder , created_At , updated_at) VALUES (?,?, ? , ? , ? , ? , ? , ? , ? , ? , ?)", ds.tableName)
	now := time.Now().In(myLocation())
	deb := req.FixedTuition + req.CourseTuition + req.ExtraOption
	remained := req.CreditAmount - req.DebitAmount
	_, err = ds.db.ExecContext(ctx, insertQuery, newID, req.StudentID, req.CourseID, req.FixedTuition, req.CourseTuition, req.ExtraOption, deb, req.CreditAmount, remained, now, now)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf("failed to insert tuition: %w", err)
	}
	return ds.selectTuitionByID(ctx, newID)

}

func (ds *TuitionDBDS) selectTuitionByID(ctx context.Context, ID int64) (res dataModels.Tuition, err error) {
	var tuition dataModels.Tuition
	readQuery := fmt.Sprintf(`
        SELECT row, student_id,course_id, fixed_tuition, course_tuition, extra_option, 	debit_amount ,credit_amount , reminder, created_at, updated_at , deleted_at
        FROM %s
        WHERE id = ? `, ds.tableName)
	var createdAt, updatedAt, deletedAt sql.NullTime
	err = ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&tuition.Row, &tuition.StudentID, &tuition.CourseID, &tuition.FixedTuition, &tuition.CourseTuition, &tuition.ExtraOption, &tuition.DebitAmount, &tuition.CreditAmount, &tuition.Reminder, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return dataModels.Tuition{}, fmt.Errorf(err.Error())
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
