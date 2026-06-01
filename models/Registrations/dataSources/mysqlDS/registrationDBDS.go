package mysqlDS

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
	"MyProject/models/Registrations/dataSources"
	"MyProject/pkg/pagination"
	TimeLoc "MyProject/pkg/timeLoc"
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

func NewEnrollmentDBDS(tableName string, db *sql.DB) (dataSources.RegistrationDS, error) {
	ff := &RegistrationDBDS{
		tableName: tableName,
		db:        db,
	}

	return ff, nil

}

func (ds *RegistrationDBDS) RegistrationsStudent(ctx context.Context, req registrationSchema.RegisterStudentRequest, role string, ID int64) (res []dataModels.ListSelectOfferingResponse, err error) {
	var studentID int64
	var register []dataModels.ListSelectOfferingResponse
	if role == "student" {
		studentID = ID
	}
	if role != "student" {
		studentID = req.StudentID
	}
	var checkStudent bool
	studentQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM student WHERE id = ? AND deleted_at IS NULL) THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, studentQuery, studentID).Scan(&checkStudent)
	if err != nil {
		return nil, err
	}
	if !checkStudent {
		return nil, errors.New("this student doesn't exist")
	}
	for _, sel := range req.Selection {
		err = ds.validateOfferingMatch(ctx, sel.OfferingID, sel.CourseNumber)
		if err != nil {
			register = append(register, dataModels.ListSelectOfferingResponse{OfferingRow: sel.OfferingID, CourseNumber: sel.CourseNumber, Err: err.Error()})
			continue
		}
		sin, err := ds.registerSingleCourse(ctx, studentID, sel.OfferingID, sel.CourseNumber)
		if err != nil {
			register = append(register, dataModels.ListSelectOfferingResponse{OfferingRow: sel.OfferingID, CourseNumber: sel.CourseNumber, Err: err.Error()})

			continue
		}

		register = append(register, dataModels.ListSelectOfferingResponse{Res: sin})
	}
	return register, nil
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
	now := time.Now().In(TimeLoc.MyLocation())
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
	now := time.Now().In(TimeLoc.MyLocation())
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
	selectQuery := fmt.Sprintf("SELECT ID, student_id,course_number, offering_row, status, enrolled_at, canceled_at, created_at, updated_at , deleted_at FROM %s LIMIT ? OFFSET ?", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var register dataModels.Registration
		err = rows.Scan(&register.ID, &register.StudentID, &register.CourseNumber, &register.OfferingRow, &register.Status, &register.EnrolledAt, &register.CanceledAt, &register.CreatedAt, &register.UpdatedAt, &register.DeletedAt)
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
	now := time.Now().In(TimeLoc.MyLocation())
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
func (ds *RegistrationDBDS) ListOfferingsStudent(ctx context.Context, req registrationSchema.ListOfferingRequest) (res []dataModels.Student, total int, err error) {
	var registers []dataModels.Student
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE offering_row = ? ", ds.tableName) // اگر میخواهید تعداد کل برای آن دانشجو را بگیرید
	err = ds.db.QueryRowContext(ctx, countQuery, req.OfferingRow).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf(err.Error(), "ERROR")
	}
	selectQuery := fmt.Sprintf("SELECT student_id , status FROM %s WHERE offering_row = ? ORDER BY student_id LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, req.OfferingRow, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var register dataModels.Student
		err = rows.Scan(&register.StudentID, &register.Status)

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
        SELECT ID, student_id,course_number, offering_row, status, enrolled_at, canceled_at, created_at, updated_at , deleted_at
        FROM %s
        WHERE ID = ? `, ds.tableName)
	err := ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&register.ID, &register.StudentID, &register.CourseNumber, &register.OfferingRow, &register.Status, &register.EnrolledAt, &register.CanceledAt, &register.CreatedAt, &register.UpdatedAt, &register.DeletedAt)
	if err != nil {
		return dataModels.Registration{}, fmt.Errorf(err.Error())
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

func (ds *RegistrationDBDS) validateOfferingMatch(ctx context.Context, offeringID, courseNumber int64) error {

	var checkOffering bool
	offeringQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND course_number = ? AND isActive = true AND capacity > 0 ) THEN 1 ELSE 0 END`
	err := ds.db.QueryRowContext(ctx, offeringQuery, offeringID, courseNumber).Scan(&checkOffering)
	if err != nil {
		return errors.New("checkOffering error")
	}
	if !checkOffering {
		return errors.New("this active offering doesn't exist or this is deActive")
	}
	return nil
}

func (ds *RegistrationDBDS) registerSingleCourse(ctx context.Context, ID int64, offering int64, courseNumber int64) (res dataModels.Registration, err error) {
	var add int64
	now := time.Now().In(TimeLoc.MyLocation())
	var alreadyRegistered bool
	checkDuplicateQuery := `
SELECT EXISTS (SELECT 1 FROM registration 
WHERE student_id = ? AND offering_row = ? AND deleted_at IS NULL)
 `
	err = ds.db.QueryRowContext(ctx, checkDuplicateQuery, ID, offering).Scan(&alreadyRegistered)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if alreadyRegistered {
		return dataModels.Registration{}, errors.New("student already registered for this course")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (student_id , course_number, offering_row,status, enrolled_at, created_at, updated_at , deleted_at) VALUES (?,?,?, ?, ?, ?, ? , ?)", ds.tableName)
	var checkCapacity bool
	studentQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND capacity > enrolled_count ) THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, studentQuery, offering).Scan(&checkCapacity)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if !checkCapacity {
		lockQuery := `SELECT row FROM offerings WHERE row = ? FOR UPDATE`
		_, err = ds.db.ExecContext(ctx, lockQuery, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		var reserved = constants.StatusReserveation
		reserve := fmt.Sprintf("UPDATE offerings SET reserveation = reserveation + 1  WHERE row = ?")
		_, err = ds.db.ExecContext(ctx, reserve, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		result, err := ds.db.ExecContext(ctx, insertQuery, ID, courseNumber, offering, reserved, now, now, now, nil)
		if err != nil {
			return dataModels.Registration{}, errors.New("you can't reserve the reservation")
		}
		add, err = result.LastInsertId()
		if err != nil {
			return dataModels.Registration{}, err
		}

	} else {
		var enrolled = constants.StatusEnrolled
		lockQuery := `SELECT row FROM offerings WHERE row = ? FOR UPDATE`
		_, err = ds.db.ExecContext(ctx, lockQuery, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		enroll := fmt.Sprintf("UPDATE offerings SET enrolled_count = enrolled_count + 1 WHERE row = ?")
		_, err = ds.db.ExecContext(ctx, enroll, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		sdd, err := ds.db.ExecContext(ctx, insertQuery, ID, courseNumber, offering, enrolled, now, now, now, nil)
		if err != nil {
			return dataModels.Registration{}, fmt.Errorf("you can't enroll the student", err)
		}
		add, err = sdd.LastInsertId()
		if err != nil {
			return dataModels.Registration{}, err
		}
	}
	return ds.readQuery(ctx, add)
}
