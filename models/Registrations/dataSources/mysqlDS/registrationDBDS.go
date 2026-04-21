package mysqlDS

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
	"MyProject/models/Registrations/dataSources"
	"MyProject/pkg/pagination"
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
			panic(r)
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
		result, err := tx.ExecContext(ctx, insertQuery, req.StudentID, req.OfferingID, reserved, now, now, now, nil)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't reserve the reservation")
		}
		add, err = result.LastInsertId()
		if err != nil {
			return dataModels.Registration{}, err
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
	if err = tx.Commit(); err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, add)
}

func (ds *RegistrationDBDS) GetRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error) {
	err = ds.check(ctx, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, req.ID)
}

func (ds *RegistrationDBDS) UpdateRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error) {
	err = ds.check(ctx, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	now := time.Now().In(myLocation())
	updateQuery := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE ID = ? ", ds.tableName)
	result, err := ds.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return dataModels.Registration{}, err
	}
	defer result.Close()
	_, err = result.ExecContext(ctx, now, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, req.ID)
}

func (ds *RegistrationDBDS) DeleteRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error) {
	err = ds.check(ctx, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	now := time.Now().In(myLocation())
	deleteQuery := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE ID = ? AND deleted_at IS NULL ", ds.tableName)
	result, err := ds.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return dataModels.Registration{}, err
	}
	defer result.Close()
	_, err = result.ExecContext(ctx, now, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, req.ID)

}

func (ds *RegistrationDBDS) ListAllRegisterStudent(ctx context.Context, req registrationSchema.SelectPageRegisteredStudentsRequest) (res []dataModels.Registration, total int, err error) {
	var registers []dataModels.Registration
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, errors.New("error getting the total count")
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var register dataModels.Registration
		var createAt, updatedAt, deletedAt, canceledAt sql.NullTime
		err = rows.Scan(&register.ID, &register.StudentID, &register.OfferingRow, &register.Status, &register.EnrolledAt, &canceledAt, &createAt, &updatedAt, &deletedAt)
		if createAt.Valid {
			register.CreatedAt = createAt.Time
		}
		if updatedAt.Valid {
			register.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			register.DeletedAt = deletedAt.Time
		}
		if canceledAt.Valid {
			register.CanceledAt = canceledAt.Time
		}
		if err != nil {
			return nil, 0, errors.New("error scanning the row")
		}
		registers = append(registers, register)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	return registers, totalRows, nil
}

func (ds *RegistrationDBDS) CancelRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest) (res dataModels.Registration, err error) {
	now := time.Now().In(myLocation())
	tx, err := ds.db.BeginTx(ctx, nil)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()
	var checkStatus bool
	var can = constants.StatusCanceled
	checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE id = ? AND status != ? AND deleted_at IS NULL) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, checkQuery, req.ID, can).Scan(&checkStatus)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if !checkStatus {
		return dataModels.Registration{}, errors.New(" status is canceled or registration deleted")
	}

	selectStatus := fmt.Sprintf("SELECT status FROM registration WHERE id = ? AND deleted_at IS NULL")
	err = tx.QueryRowContext(ctx, selectStatus, req.ID).Scan(&res.Status)
	var canceling = constants.StatusCanceled
	updateQuery := fmt.Sprintf("UPDATE %s SET canceled_at = ? , status = ? WHERE ID = ? AND status = ?", ds.tableName)
	result, err := tx.PrepareContext(ctx, updateQuery)
	if err != nil {
		return dataModels.Registration{}, err
	}
	defer result.Close()
	_, err = result.ExecContext(ctx, now, canceling, req.ID, res.Status)
	if err != nil {
		return dataModels.Registration{}, err
	}
	selectOfferingRow := fmt.Sprintf("SELECT offering_row FROM %s WHERE id = ? ", ds.tableName)
	var offeringRow int
	err = tx.QueryRowContext(ctx, selectOfferingRow, req.ID).Scan(&offeringRow)
	if err != nil {
		return dataModels.Registration{}, fmt.Errorf("cannot find offering row for registration %d: %w", req.ID, err)
	}

	if res.Status == constants.StatusReserveation {

		decrementEnrolledQuery := fmt.Sprintf("UPDATE offerings SET reserveation = reserveation - 1  WHERE row = ? AND reserveation > 0")
		result, err = tx.PrepareContext(ctx, decrementEnrolledQuery)
		if err != nil {
			return dataModels.Registration{}, err
		}
		defer result.Close()
		_, err = result.ExecContext(ctx, offeringRow)
		if err != nil {
			return dataModels.Registration{}, err
		}
	}
	if res.Status == constants.StatusEnrolled {

		decrementEnrolledQuery := fmt.Sprintf("UPDATE offerings SET enrolled_count = enrolled_count - 1 WHERE row = ? AND enrolled_count > 0")
		result, err = tx.PrepareContext(ctx, decrementEnrolledQuery)
		if err != nil {
			return dataModels.Registration{}, err
		}
		defer result.Close()
		_, err = result.ExecContext(ctx, offeringRow)
		if err != nil {
			return dataModels.Registration{}, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, req.ID)

}

func (ds *RegistrationDBDS) ListStudentsOffering(ctx context.Context, req registrationSchema.ListStudentsRequest) (res []dataModels.Offering, total int, err error) {
	var registers []dataModels.Offering
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE student_id = ? ", ds.tableName) // اگر میخواهید تعداد کل برای آن دانشجو را بگیرید
	err = ds.db.QueryRowContext(ctx, countQuery, req.StudentID).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf(err.Error(), "ERROR")
	}
	selectQuery := fmt.Sprintf("SELECT offering_row , status FROM %s WHERE student_id = ? ORDER BY offering_row LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, req.StudentID, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var register dataModels.Offering
		err = rows.Scan(&register.OfferingRow, &register.Status)

		if err != nil {
			return nil, 0, errors.New("error scanning the row")
		}

		registers = append(registers, register)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	return registers, totalRows, nil
}

func (ds *RegistrationDBDS) readQuery(ctx context.Context, ID int64) (dataModels.Registration, error) {
	var register dataModels.Registration
	readQuery := fmt.Sprintf(`
        SELECT ID, student_id, offering_row, status, enrolled_at, canceled_at, created_at, updated_at , deleted_at
        FROM %s
        WHERE id = ? `, ds.tableName)
	var createdAt, updatedAt, deletedAt, canceledAt sql.NullTime
	err := ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&register.ID, &register.StudentID, &register.OfferingRow, &register.Status, &register.EnrolledAt, &canceledAt, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return dataModels.Registration{}, fmt.Errorf(err.Error())
	}

	if canceledAt.Valid {
		register.CanceledAt = canceledAt.Time.In(myLocation())
	}
	if createdAt.Valid {
		register.CreatedAt = createdAt.Time.In(myLocation())
	} else {
		register.CreatedAt = time.Time{}
	}

	if updatedAt.Valid {
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

func (ds *RegistrationDBDS) check(ctx context.Context, id int64) error {
	var checkRegister bool
	selectQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, selectQuery, id).Scan(&checkRegister)
	if err != nil {
		return err
	}
	if !checkRegister {
		return errors.New("you can't check the registration . because there is no registration")
	}
	return nil

}
