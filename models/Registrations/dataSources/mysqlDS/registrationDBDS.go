package mysqlDS

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
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

func NewEnrollmentDBDS(tableName string, db *sql.DB) (dataModels.Registration, error) {
	ff := &RegistrationDBDS{
		tableName: tableName,
		db:        db,
	}

	return ff, nil

}

func (ds *RegistrationDBDS) RegistrationsStudent(ctx context.Context, req registrationSchema.RegisterStudentRequest) (res dataModels.Registration, err error) {
	var registration dataModels.Registration
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
CASE WHEN EXISTS (SELECT 1 FROM registrations WHERE student_id = ?)`
	err = tx.QueryRow(teacherQuery, req.StudentID).Scan(&checkStudent)
	if err != nil {
		return dataModels.Registration{}, errors.New("checkStudent error")
	}
	if !checkStudent {
		return dataModels.Registration{}, errors.New("this student doesn't exist")
	}
	var checkOffering bool
	teacherQuery = `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registrations WHERE offering_row = ? AND isActive = true )`
	err = tx.QueryRow(teacherQuery, req.OfferingID).Scan(&checkOffering)
	if err != nil {
		return dataModels.Registration{}, errors.New("checkOffering error")
	}
	if !checkOffering {
		return dataModels.Registration{}, errors.New("this active offering doesn't exist")
	}

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
