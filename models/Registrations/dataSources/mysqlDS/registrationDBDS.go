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
	var add int64
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
CASE WHEN EXISTS (SELECT 1 FROM student WHERE id = ? AND deleted_at IS NULL) THEN 1 ELSE 0 END`
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
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND isActive = true AND capacity > 0  ) THEN 1 ELSE 0 END`
	err = tx.QueryRow(teacherQuery, req.OfferingID).Scan(&checkOffering)
	if err != nil {
		return dataModels.Registration{}, errors.New("checkOffering error")
	}
	if !checkOffering {
		return dataModels.Registration{}, errors.New("this active offering doesn't exist or this is deActive")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (student_id, offering_row,status, enrolled_at, created_at, updated_at , deleted_at) VALUES (?,?, ?, ?, ?, ? , ?)", ds.tableName)
	var checkCapacity bool
	teacherQuery = `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND capacity > enrolled_count ) THEN 1 ELSE 0 END`
	err = tx.QueryRow(teacherQuery, req.OfferingID).Scan(&checkCapacity)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if !checkCapacity {
		var reserved = constants.StatusReserveation
		reserve := fmt.Sprintf("UPDATE offerings SET reserveation = reserveation + 1  WHERE row = ?")
		_, err = tx.Exec(reserve, req.OfferingID)
		if err != nil {
			return dataModels.Registration{}, err
		}
		_, err = tx.ExecContext(ctx, insertQuery, req.StudentID, req.OfferingID, reserved, now, now, now, nil)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't reserve the reservation")
		}

	} else {
		var enrolled = constants.StatusEnrolled
		enroll := fmt.Sprintf("UPDATE offerings SET enrolled_count = enrolled_count + 1 WHERE row = ?")
		_, err = tx.Exec(enroll, req.OfferingID)
		if err != nil {
			return dataModels.Registration{}, err
		}
		sdd, err := tx.ExecContext(ctx, insertQuery, req.StudentID, req.OfferingID, enrolled, now, now, now, nil)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't enroll the student")
		}
		add, err = sdd.LastInsertId()
		if err != nil {
			return dataModels.Registration{}, err
		}
	}
	return ds.readQuery(ctx, add)
}

func (ds *RegistrationDBDS) readQuery(ctx context.Context, ID int64) (dataModels.Registration, error) {
	var register dataModels.Registration
	readQuery := fmt.Sprintf(`
        SELECT ID, student_id, offering_row, status, enrolled_at, canceled_at, created_at, updated_at , deleted_at
        FROM %s
        WHERE id = ? `, ds.tableName)
	var createdAt, updatedAt, deletedAt sql.NullTime
	err := ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&register.ID, &register.StudentID, &register.OfferingRow, &register.Status, &register.EnrolledAt, &register.CanceledAt, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return dataModels.Registration{}, fmt.Errorf(err.Error(), ID, err, ID, err)
	}
	if createdAt.Valid {
		register.CreatedAt = createdAt.Time.In(myLocation())
	} else {
		register.CreatedAt = time.Time{}
	}

	if updatedAt.Valid {
		fmt.Println(updatedAt.Time)
		register.UpdatedAt = updatedAt.Time.In(myLocation())
	} else {
		register.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		register.DeletedAt = deletedAt.Time.In(myLocation())
	} else {
		register.DeletedAt = time.Time{}
	}

	return register, nil

}
