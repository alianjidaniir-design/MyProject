package mysqlDS

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
	"MyProject/models/Registrations/dataSources"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type RegistrationDBDS struct {
	tableName string
	db        *sql.DB
}

func myLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return location
}

func NewEnrollmentDBDS(tableName string, db *sql.DB) (dataSources.RegistrationDS, error) {
	ff := &RegistrationDBDS{
		tableName: tableName,
		db:        db,
	}

	return ff, nil

}

func (ds *RegistrationDBDS) RegistrationsStudent(ctx context.Context, req registrationSchema.RegisterStudentRequest) (res dataModels.Registration, err error) {
	now := time.Now().In(myLocation())
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	var checkStudent bool
	teacherQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM student WHERE id = ?)`
	err = tx.QueryRow(teacherQuery, req.StudentID).Scan(&checkStudent)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if !checkStudent {
		return dataModels.Registration{}, errors.New("this student doesn't exist")
	}
	var checkOffering bool
	teacherQuery = `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND isActive = true AND capacity > 0  )`
	err = tx.QueryRow(teacherQuery, req.OfferingID).Scan(&checkOffering)
	if err != nil {
		return dataModels.Registration{}, errors.New("checkOffering error")
	}
	if !checkOffering {
		return dataModels.Registration{}, errors.New("this active offering doesn't exist")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (student_id, offering_row,status, enrolled_at, created_at, updated_at) VALUES (?,?, ?, ?, ?, ?)", ds.tableName)
	var checkCapacity bool
	teacherQuery = `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE id = ? AND capacity > enrolled_count )`
	err = tx.QueryRow(teacherQuery, req.OfferingID).Scan(&checkCapacity)
	if err != nil {
		return dataModels.Registration{}, errors.New("checkCapacity error")
	}
	if !checkCapacity {
		var reserved = constants.StatusReserveation
		reserve := fmt.Sprintf("UPDATE offerings SET reserveation = reserveation + 1  WHERE id = ?")
		_, err = tx.Exec(reserve, req.OfferingID)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't reserve the reservation")
		}
		_, err := tx.ExecContext(ctx, insertQuery, req.StudentID, req.OfferingID, reserved, now, now, now)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't reserve the reservation")
		}
		var enrolled = constants.StatusEnrolled
		enroll := fmt.Sprintf("UPDATE offerings SET enrolled = enrolled + 1 WHERE id = ?")
		_, err = tx.Exec(enroll, req.OfferingID)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't reserve the reservation")
		}
		_, err = tx.ExecContext(ctx, insertQuery, req.StudentID, req.OfferingID, enrolled, now, now, now)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't enroll the student")
		}
	}
	var lastID int64
	selectQuery := fmt.Sprintf("SELECT id FROM %s ", ds.tableName)
	err = tx.QueryRow(selectQuery).Scan(&lastID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, lastID)
}

func (ds *RegistrationDBDS) readQuery(ctx context.Context, ID int64) (dataModels.Registration, error) {
	var enrollment dataModels.Registration
	readQuery := fmt.Sprintf(`
        SELECT id, student_id, offering_row, status, enrolled_at, canceled_at, created_at, updated_at , deleted_at
        FROM %s
        WHERE id = ? AND status = ? `, ds.tableName)
	var ff = constants.StatusEnrolled
	err := ds.db.QueryRowContext(ctx, readQuery, ID, ff).Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.OfferingRow, &enrollment.Status, &enrollment.EnrolledAt, &enrollment.CanceledAt, &enrollment.CreatedAt, &enrollment.UpdatedAt, &enrollment.DeletedAt)
	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return dataModels.Registration{}, err
	}

	return enrollment, nil

}
