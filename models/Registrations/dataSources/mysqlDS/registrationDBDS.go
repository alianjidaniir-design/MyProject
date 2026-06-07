package mysqlDS

import (
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataModels"
	"MyProject/models/Registrations/dataSources"
	"MyProject/models/student/dataModel"
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

func NewRegisterDBDS(tableName string, db *sql.DB) (dataSources.RegistrationDS, error) {
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
	} else {
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

	total, err := ds.checkingTotalUnits(ctx, studentID)
	if err != nil {
		return nil, err
	}

	maxim, err := ds.maxNumberUnits(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var addedUnits int64 = 0

	for _, sel := range req.Selection {
		courseUnits, err := ds.GetLastCourseUnits(ctx, sel.OfferingID)
		if err != nil {
			register = append(register, dataModels.ListSelectOfferingResponse{
				StudentID:    studentID,
				OfferingRow:  sel.OfferingID,
				CourseNumber: sel.CourseNumber,
				Status:       "",
				Err:          fmt.Sprintf("offering not found: %v", err),
			})
			continue
		}

		err = ds.validateOfferingMatch(ctx, sel.OfferingID, sel.CourseNumber)
		if err != nil {
			register = append(register, dataModels.ListSelectOfferingResponse{StudentID: studentID, OfferingRow: sel.OfferingID, CourseNumber: sel.CourseNumber, Status: "", Err: err.Error()})
			continue
		}
		if total+addedUnits+courseUnits > maxim {
			register = append(register, dataModels.ListSelectOfferingResponse{StudentID: studentID, OfferingRow: sel.OfferingID, CourseNumber: sel.CourseNumber, Status: "", Err: errors.New("selects is more than maxim units").Error()})
			continue
		}
		sin, err := ds.registerSingleCourse(ctx, studentID, sel.OfferingID, sel.CourseNumber, role, sel.IsReserve)
		if err != nil {
			register = append(register, dataModels.ListSelectOfferingResponse{StudentID: studentID, OfferingRow: sel.OfferingID, CourseNumber: sel.CourseNumber, Status: sin.Status, Err: err.Error()})
			continue
		}

		addedUnits += courseUnits

		register = append(register, dataModels.ListSelectOfferingResponse{ID: sin.ID, StudentID: studentID, OfferingRow: sel.OfferingID, CourseNumber: sel.CourseNumber, Status: sin.Status})
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

func (ds *RegistrationDBDS) DeleteRegisterStudent(ctx context.Context, req registrationSchema.GetRegisteredStudentsRequest, role string, studentID int64) (res dataModels.Registration, err error) {
	err = ds.check(ctx, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if role == "student" {
		err = ds.isCorrect(ctx, req.ID, studentID)
		if err != nil {
			return dataModels.Registration{}, err
		}
		reg, err := ds.readQuery(ctx, req.ID)
		if err != nil {
			return dataModels.Registration{}, err
		}
		if reg.Registrar != "student" {
			return dataModels.Registration{}, errors.New("this is registered by admin . you can not deleted it")
		}
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE ID = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return dataModels.Registration{}, err
	}
	return dataModels.Registration{}, nil

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
	selectQuery := fmt.Sprintf("SELECT ID, student_id,course_number, offering_row, status, enrolled_at, canceled_at, created_at, updated_at FROM %s LIMIT ? OFFSET ?", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var register dataModels.Registration
		err = rows.Scan(&register.ID, &register.StudentID, &register.CourseNumber, &register.OfferingRow, &register.Status, &register.EnrolledAt, &register.CanceledAt, &register.CreatedAt, &register.UpdatedAt)
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
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE id = ? AND status != ?) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, checkQuery, req.ID, can).Scan(&checkStatus)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if !checkStatus {
		return dataModels.Registration{}, errors.New(" status is canceled or registration deleted")
	}

	selectStatus := fmt.Sprintf("SELECT status FROM registration WHERE id = ?")
	err = tx.QueryRowContext(ctx, selectStatus, req.ID).Scan(&res.Status)
	var canceling = constants.StatusCanceled
	updateQuery := fmt.Sprintf("UPDATE %s SET canceled_at = ? , updated_at = ? , status = ? WHERE ID = ? AND status = ?", ds.tableName)
	result, err := tx.PrepareContext(ctx, updateQuery)
	if err != nil {
		return dataModels.Registration{}, err
	}
	defer result.Close()
	_, err = result.ExecContext(ctx, now, now, canceling, req.ID, res.Status)
	if err != nil {
		return dataModels.Registration{}, err
	}
	selectOfferingRow := fmt.Sprintf("SELECT offering_row FROM %s WHERE id = ? ", ds.tableName)
	var offeringRow int64
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

func (ds *RegistrationDBDS) ListClassesStudent(ctx context.Context, req registrationSchema.Pages, studentID int64) (res []dataModels.TermClassSchedules, total int, page int, err error) {
	var termClassSchedules []dataModels.TermClassSchedules
	var classes []dataModels.DetailClassScheduler
	var termID int64
	getTermIDQuery := `
        SELECT id 
        FROM terms 
        WHERE term = ? AND year = ?
    `
	err = ds.db.QueryRowContext(ctx, getTermIDQuery, req.Term, req.Year).Scan(&termID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, 0, fmt.Errorf("term %d for year %d does not exist. Available terms: check terms table", req.Term, req.Year)
		}
		return nil, 0, 0, err
	}
	page, size, err := pagination.CheckPage(req.Page, constants.PageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	limit := size
	offset := (page - 1) * limit
	var totalRows int
	countQuery := `
        SELECT COUNT(*) 
        FROM registration r
        JOIN offerings o ON r.offering_row = o.row
        WHERE r.student_id = ? 
          AND o.term_id = ?
    `
	err = ds.db.QueryRowContext(ctx, countQuery, studentID, termID).Scan(&totalRows)
	if err != nil {
		return nil, 0, 0, err
	}
	var tot int
	selectQuery := `
SELECT
o.row AS offering_row,
o.group_number AS offering_group_number,
c.course_number AS course_number,
c.title AS title,
c.unit AS unit,
t.name AS teacher_name,
t.last_name AS teacher_last_name,
o.class_start_time AS class_start_time,
o.class_end_time AS class_end_time
FROM registration r
JOIN offerings o ON r.offering_row = o.row
JOIN courses c ON o.course_number = c.course_number
JOIN teachers t ON o.teacher_id = t.ID
JOIN student u ON r.student_id = u.ID
WHERE r.student_id = ? AND o.term_id = ?
ORDER BY o.class_start_time LIMIT ? OFFSET ?;
`
	rows, err := ds.db.QueryContext(ctx, selectQuery, studentID, termID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var class dataModels.DetailClassScheduler
		err = rows.Scan(&class.OfferingRow, &class.OfferingGroupNumber, &class.CourseNumber, &class.Title, &class.Unit, &class.TeacherName, &class.TeacherLastName, &class.ClassStartTime, &class.ClassEndTime)
		if err != nil {
			return nil, 0, 0, err
		}
		sumUnits := class.Unit
		tot += sumUnits
		classes = append(classes, class)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}

	var termClassSchedulers dataModels.TermClassSchedules

	termClassSchedulers.Term = req.Term
	termClassSchedulers.Year = req.Year
	termClassSchedulers.Classes = classes
	termClassSchedulers.TotalUnits = tot
	termClassSchedules = append(termClassSchedules, termClassSchedulers)

	return termClassSchedules, totalRows, page, nil

}

func (ds *RegistrationDBDS) readQuery(ctx context.Context, ID int64) (dataModels.Registration, error) {
	var register dataModels.Registration
	readQuery := fmt.Sprintf(`
        SELECT ID, student_id,course_number, offering_row, status,registrar, enrolled_at, canceled_at, created_at, updated_at
        FROM %s
        WHERE ID = ? `, ds.tableName)
	err := ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&register.ID, &register.StudentID, &register.CourseNumber, &register.OfferingRow, &register.Status, &register.Registrar, &register.EnrolledAt, &register.CanceledAt, &register.CreatedAt, &register.UpdatedAt)
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

func (ds *RegistrationDBDS) registerSingleCourse(ctx context.Context, ID int64, offering int64, courseNumber int64, role string, reserve bool) (res dataModels.Registration, err error) {
	var add int64
	now := time.Now().In(TimeLoc.MyLocation())
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return dataModels.Registration{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()
	var alreadyRegistered bool
	checkDuplicateQuery := `
SELECT EXISTS (SELECT 1 FROM registration 
WHERE student_id = ? AND offering_row = ? )
 `
	err = tx.QueryRowContext(ctx, checkDuplicateQuery, ID, offering).Scan(&alreadyRegistered)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if alreadyRegistered {
		return dataModels.Registration{}, errors.New("student already registered or reservation for this course")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (student_id , course_number, offering_row,status,registrar, enrolled_at, created_at, updated_at ) VALUES (?,?,?, ?, ?, ?, ? , ?)", ds.tableName)
	var checkCapacity bool
	studentQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND capacity > enrolled_count ) THEN 1 ELSE 0 END`
	err = tx.QueryRowContext(ctx, studentQuery, offering).Scan(&checkCapacity)
	if err != nil {
		return dataModels.Registration{}, err
	}
	if !checkCapacity {
		if reserve != true {
			return dataModels.Registration{}, errors.New("The class is full. You can reserve it.")
		}
		lockQuery := `SELECT row FROM offerings WHERE row = ? FOR UPDATE`
		_, err = tx.ExecContext(ctx, lockQuery, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		var reserved = constants.StatusReserveation
		reserve := fmt.Sprintf("UPDATE offerings SET reserveation = reserveation + 1  WHERE row = ?")
		_, err = tx.ExecContext(ctx, reserve, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		result, err := tx.ExecContext(ctx, insertQuery, ID, courseNumber, offering, reserved, role, now, now, now)
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
		_, err = tx.ExecContext(ctx, lockQuery, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		enroll := fmt.Sprintf("UPDATE offerings SET enrolled_count = enrolled_count + 1 WHERE row = ?")
		_, err = tx.ExecContext(ctx, enroll, offering)
		if err != nil {
			return dataModels.Registration{}, err
		}
		sdd, err := tx.ExecContext(ctx, insertQuery, ID, courseNumber, offering, enrolled, role, now, now, now)
		if err != nil {
			return dataModels.Registration{}, fmt.Errorf("you can't enroll the student", err)
		}
		add, err = sdd.LastInsertId()
		if err != nil {
			return dataModels.Registration{}, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Registration{}, err
	}
	return ds.readQuery(ctx, add)
}

func (ds *RegistrationDBDS) checkingTotalUnits(ctx context.Context, studentID int64) (tot int64, err error) {
	var count int64
	query := `
        SELECT COALESCE(SUM(c.unit), 0)
        FROM registration r
        JOIN offerings o ON r.offering_row = o.row
        JOIN courses c ON o.course_number = c.course_number
        WHERE r.student_id = ? 
          AND (r.status = 'enrolled' OR r.status = 'reserveation')
          `
	err = ds.db.QueryRowContext(ctx, query, studentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ds *RegistrationDBDS) ListClassesTeacher(ctx context.Context, req registrationSchema.Pages, teacherID int64) (res []dataModels.TermClasses, total int, page int, err error) {
	var termClass []dataModels.TermClasses
	var classes []dataModels.DetailClasses
	var termID int64
	getTermIDQuery := `
        SELECT id 
        FROM terms 
        WHERE term = ? AND year = ?
    `
	err = ds.db.QueryRowContext(ctx, getTermIDQuery, req.Term, req.Year).Scan(&termID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, 0, fmt.Errorf("term %d for year %d does not exist. Available terms: check terms table", req.Term, req.Year)
		}
		return nil, 0, 0, err
	}
	page, size, err := pagination.CheckPage(req.Page, constants.PageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	limit := size
	offset := (page - 1) * limit
	var totalRows int
	countQuery := `
        SELECT COUNT(*) 
        FROM registration r
        JOIN offerings o ON r.offering_row = o.row
        WHERE o.teacher_id = ? 
          AND o.term_id = ?
    `
	err = ds.db.QueryRowContext(ctx, countQuery, teacherID, termID).Scan(&totalRows)
	if err != nil {
		return nil, 0, 0, err
	}
	selectQuery := `
SELECT
o.group_number AS offering_group_number,
c.course_number AS course_number,
c.title AS title,
o.capacity AS capacity,
o.enrolled_count AS enrolled_count,
o.class_start_time AS class_start_time,
o.class_end_time AS class_end_time,
c.department_id AS department_id
FROM registration r
JOIN offerings o ON r.offering_row = o.row
JOIN courses c ON o.course_number = c.course_number
JOIN teachers t ON o.teacher_id = t.ID
WHERE o.teacher_id = ? AND o.term_id = ?
ORDER BY o.class_start_time LIMIT ? OFFSET ?;
`
	rows, err := ds.db.QueryContext(ctx, selectQuery, teacherID, termID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var class dataModels.DetailClasses
		err = rows.Scan(&class.OfferingGroupNumber, &class.CourseNumber, &class.Title, &class.Capacity, &class.EnrolledCount, &class.ClassStartTime, &class.ClassEndTime, &class.DepartmentID)
		if err != nil {
			return nil, 0, 0, err
		}

		classes = append(classes, class)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}

	var termClassSchedulers dataModels.TermClasses

	termClassSchedulers.Term = req.Term
	termClassSchedulers.Year = req.Year
	termClassSchedulers.Classes = classes
	termClass = append(termClass, termClassSchedulers)

	return termClass, totalRows, page, nil

}

func (ds *RegistrationDBDS) GetLastCourseUnits(ctx context.Context, offeringID int64) (int64, error) {
	var units int64
	query := `
        SELECT COALESCE(c.unit, 0)
        FROM offerings o
        JOIN courses c ON o.course_number = c.course_number
        WHERE o.row = ?
    `
	err := ds.db.QueryRowContext(ctx, query, offeringID).Scan(&units)
	return units, err
}

func (ds RegistrationDBDS) maxNumberUnits(ctx context.Context, StudentID int64) (int64, error) {
	var student dataModel.Student
	checkStudent := fmt.Sprintf("SELECT level FROM student WHERE ID = ?")
	err := ds.db.QueryRowContext(ctx, checkStudent, StudentID).Scan(&student.Level)
	if err != nil {
		return 0, err
	}
	switch student.Level {
	case "bachelor":
		return 24, nil
	case "master":
		return 16, nil
	case "phd":
		return 12, nil
	default:
		return 0, errors.New("invalid student level")
	}

}

func (ds *RegistrationDBDS) isCorrect(ctx context.Context, ID int64, studentID int64) error {
	var checkRegister bool
	selectQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE ID = ? AND student_id = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, selectQuery, ID, studentID).Scan(&checkRegister)
	if err != nil {
		return err
	}
	if !checkRegister {
		return errors.New("there is not registration")
	}
	return nil

}
