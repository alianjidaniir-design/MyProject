package mySqlDS

import (
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering/dataModels"
	"MyProject/pkg/filter"
	"MyProject/pkg/pagination"
	"MyProject/pkg/timeLoc"
	"MyProject/pkg/val"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type OfferingDBDS struct {
	tableName string
	db        *sql.DB
}

func NewOfferingDBDS(tableName string, db *sql.DB) (*OfferingDBDS, error) {
	offer := &OfferingDBDS{
		tableName: tableName,
		db:        db,
	}
	return offer, nil
}

func (ds *OfferingDBDS) CreateOffering(ctx context.Context, req offeringSchema.CreateOfferingRequest) (res dataModels.Offering, err error) {
	err = val.CheckValidation(req)
	if err != nil {
		return dataModels.Offering{}, err
	}
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	var checkCourse bool
	courseQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM courses WHERE course_number = ? ) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, courseQuery, req.CourseNumber).Scan(&checkCourse)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !checkCourse {
		return dataModels.Offering{}, errors.New("course does not exist")
	}

	var checkingTeacher bool

	teacherQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM teachers WHERE id = ?) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, teacherQuery, req.TeacherId).Scan(&checkingTeacher)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !checkingTeacher {
		return res, errors.New("Teacher does not exist")
	}

	var checkTerm bool
	termQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM terms WHERE id = ?) THEN 1 ELSE 0 END`
	err = tx.QueryRowContext(ctx, termQuery, req.TermId).Scan(&checkTerm)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !checkTerm {
		return res, errors.New("Terms does not exist")
	}

	var lastID int64

	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row), 0) FROM %s", ds.tableName)
	err = tx.QueryRowContext(ctx, lastIDQuery).Scan(&lastID)
	if err != nil {
		return dataModels.Offering{}, err
	}
	var check int
	timeStart := req.ClassStartTime
	timeEnd := req.ClassEndTime
	err = timeLoc.CheckDuration(timeStart, timeEnd)
	if err != nil {
		return dataModels.Offering{}, err
	}

	start, err := timeLoc.FormatTime(timeStart)
	if err != nil {
		return dataModels.Offering{}, err
	}
	end, err := timeLoc.FormatTime(timeEnd)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if start > end {
		return dataModels.Offering{}, errors.New("start time is greater than end time")
	}

	checkUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE group_number = ?", ds.tableName)
	err = tx.QueryRowContext(ctx, checkUnique, req.GroupNumber).Scan(&check)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if check > 0 {
		return dataModels.Offering{}, errors.New("groupNumber already exists")
	}

	var count int
	if req.Week == constants.Even {
		checkingUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id = ? AND (week = ? OR week = ?) AND class_start_time < ? AND class_end_time > ?", ds.tableName)
		err = tx.QueryRowContext(ctx, checkingUnique, req.TeacherId, constants.Even, constants.All, start, end).Scan(&count)
		if err != nil {

			return dataModels.Offering{}, err
		}
		if count > 0 {
			return dataModels.Offering{}, errors.New("groupNumber already exists or It conflicts with this teacher's other class.")
		}
	} else if req.Week == constants.Odd {
		checkingUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id = ? AND (week = ? OR week = ?) AND day = ? AND class_start_time <= ? AND class_end_time >= ?", ds.tableName)
		err = tx.QueryRowContext(ctx, checkingUnique, req.TeacherId, constants.Odd, constants.All, req.Day, start, end).Scan(&count)
		if err != nil {
			return dataModels.Offering{}, err
		}
		if count > 0 {
			return dataModels.Offering{}, errors.New("groupNumber already exists or It conflicts with this teacher's other class.")
		}
	} else if req.Week == constants.All {
		checkingUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id = ? AND day = ? AND class_start_time < ? AND class_end_time > ?", ds.tableName)
		err = tx.QueryRowContext(ctx, checkingUnique, req.TeacherId, req.Day, start, end).Scan(&count)
		if err != nil {
			return dataModels.Offering{}, err
		}
		if count > 0 {
			return dataModels.Offering{}, errors.New("groupNumber already exists or It conflicts with this teacher's other class.")
		}
	} else {
		return dataModels.Offering{}, errors.New("invalid week")
	}

	newID := lastID + 1
	now := time.Now().In(timeLoc.MyLocation())
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , group_number , course_number , teacher_id , capacity , isActive ,term_id,week, day , class_start_time , class_end_time, exam_start_time, exam_finish_time , updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", ds.tableName)
	_, err = tx.ExecContext(ctx, insertQuery, newID, req.GroupNumber, req.CourseNumber, req.TeacherId, req.Capacity, req.IsActive, req.TermId, req.Week, req.Day, start, end, req.ExamStartTime, req.ExamEndTime, now)
	if err != nil {
		return dataModels.Offering{}, err
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Offering{}, err
	}
	return ds.readOfferingByID(ctx, newID)

}

func (ds *OfferingDBDS) ListOffering(ctx context.Context, req offeringSchema.ListOfferingsRequest) (res []dataModels.ListOfferings, total int, err error) {
	var offerings []dataModels.ListOfferings
	var totalRows int
	err = val.CheckValidation(req)
	if err != nil {
		return nil, 0, err
	}
	fil := []filter.Filter{

		{Con: "college", Value: req.College},
		{Con: "educational_group", Value: req.EducationalGroup},
		{Con: "week", Value: req.Week},
		{Con: "day", Value: req.Day},
	}
	cond, args := filter.Filtering(fil...)
	whereClause := "term = ? AND year = ?"
	queryArgs := []interface{}{req.Term, req.Year}

	if cond != "" && cond != "1=1" {
		whereClause += " AND " + cond
		queryArgs = append(queryArgs, args...)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM offering_list WHERE %s", whereClause)
	err = ds.db.QueryRowContext(ctx, countQuery, queryArgs...).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf("error in rows count: %w", err)
	}
	selectQuery := fmt.Sprintf("SELECT * FROM offering_list WHERE %s ORDER BY year DESC, term DESC ", whereClause)
	rows, err := ds.db.QueryContext(ctx, selectQuery, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var offering dataModels.ListOfferings
		err = rows.Scan(&offering.Row, &offering.GroupNumber, &offering.CourseNumber, &offering.Title, &offering.Unit, &offering.TeacherName, &offering.TeacherLastName, &offering.College, &offering.EducationalGroup, &offering.Capacity, &offering.EnrolledCount, &offering.IsActive, &offering.Reservation, &offering.Week, &offering.Day, &offering.ClassStartTime, &offering.ClassEndTime, &offering.ExamStartTime, &offering.ExamEndTime, &offering.Term, &offering.Year)
		if err != nil {
			return nil, 0, fmt.Errorf("Error scanning row", err.Error())
		}
		offerings = append(offerings, offering)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("Error : ", err.Error())
	}
	return offerings, totalRows, nil
}

func (ds *OfferingDBDS) GetOffering(ctx context.Context, req offeringSchema.GetRowOfferingRequest) (res dataModels.Offering, err error) {
	err = ds.checkID(ctx, req.Row)
	if err != nil {
		return res, err
	}
	return ds.readOfferingByID(ctx, req.Row)
}
func (ds *OfferingDBDS) DeActiveOffering(ctx context.Context, req offeringSchema.GetRowOfferingRequest) (res dataModels.Offering, err error) {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND isActive = true ) THEN 1 ELSE 0 END
`
	err = ds.db.QueryRowContext(ctx, search, req.Row).Scan(&check)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !check {
		return dataModels.Offering{}, errors.New("active Offering does not exist")
	}
	deActiveQuery := fmt.Sprintf("UPDATE `%s` SET isActive = 0 WHERE row = ?", ds.tableName)
	update, err := ds.db.PrepareContext(ctx, deActiveQuery)
	if err != nil {
		return dataModels.Offering{}, err
	}
	defer update.Close()
	result, err := update.ExecContext(ctx, req.Row)
	if err != nil {
		return dataModels.Offering{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return dataModels.Offering{}, err
	}
	return ds.readOfferingByID(ctx, req.Row)
}

func (ds *OfferingDBDS) EditOffering(ctx context.Context, req offeringSchema.EditOffering) (res dataModels.Offering, err error) {
	err = val.CheckValidation(req)
	if err != nil {
		return res, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	err = ds.checkID(ctx, req.Row)
	if err != nil {
		return dataModels.Offering{}, err
	}
	var enroll, reserve, capacity int
	var startExam, finishExam *time.Time
	lockQuery := `SELECT enrolled_count , reserveation , capacity , exam_start_time , exam_finish_time FROM offerings WHERE row = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, lockQuery, req.Row).Scan(&enroll, &reserve, &capacity, &startExam, &finishExam)
	if err != nil {
		return dataModels.Offering{}, err
	}

	now := time.Now().In(timeLoc.MyLocation())
	updateQuery := "UPDATE offerings SET updated_at = ? "
	args := []interface{}{now}
	if req.Capacity != 0 {
		deff := req.Capacity - capacity
		if req.Capacity < enroll {
			return dataModels.Offering{}, errors.New("capacity out of range")
		} else if reserve > 0 && deff > 0 {
			updateRegister := `
UPDATE registration
SET 
    updated_at = ? ,
    queuePosition = CASE 
    WHEN queuePosition >= ? THEN queuePosition - ?
    ELSE 0
    END
WHERE offering_row = ?
  AND status = ?
  AND queuePosition > 0;
`
			result, err := tx.ExecContext(ctx, updateRegister, now, deff, deff, req.Row, constants.Reservation)
			if err != nil {
				return dataModels.Offering{}, err
			}

			rowsAffected, _ := result.RowsAffected()
			if rowsAffected > 0 {
				updateStatus := `
					UPDATE registration
					SET status = ?
					WHERE offering_row = ? AND queuePosition = 0 AND status = ?
				`
				_, err = tx.ExecContext(ctx, updateStatus, constants.Enrolled, req.Row, constants.Reservation)
				if err != nil {
					return dataModels.Offering{}, err
				}
			}

			updateOffering := "UPDATE offerings SET enrolled_count = enrolled_count + ?, reserveation = reserveation - ? WHERE row = ?"
			if deff <= reserve {
				_, err = tx.ExecContext(ctx, updateOffering, deff, deff, req.Row)
				if err != nil {
					return dataModels.Offering{}, err
				}
			} else {
				_, err = tx.ExecContext(ctx, updateOffering, reserve, reserve, req.Row)
				if err != nil {
					return dataModels.Offering{}, err
				}
			}
		}
		updateQuery += " , capacity = ?"
		args = append(args, req.Capacity)
	}

	if req.IsActive != nil {
		if *req.IsActive == false && enroll > 0 {
			return dataModels.Offering{}, errors.New("enrolled student in course . can not deActive it")
		}
		updateQuery += ", isActive = ?"
		args = append(args, req.IsActive)
	}

	if req.ExamStartTime != "" {

		start, err := timeLoc.FormatDataTime(req.ExamStartTime)
		if err != nil {
			return dataModels.Offering{}, err
		}
		if finishExam != nil {
			err = timeLoc.CheckTimeExam(start, finishExam)
			if err != nil {
				return dataModels.Offering{}, err
			}
		}

		updateQuery += ", exam_start_time = ? "
		args = append(args, start)
	}
	if req.ExamFinishTime != "" {
		finish, err := timeLoc.FormatDataTime(req.ExamFinishTime)
		if err != nil {
			return dataModels.Offering{}, err
		}
		if startExam != nil {
			err = timeLoc.CheckTimeExam(startExam, finish)
			if err != nil {
				return dataModels.Offering{}, err
			}
		}

		updateQuery += ", exam_finish_time = ? "
		args = append(args, finish)
	}

	updateQuery += " WHERE row = ?"
	args = append(args, req.Row)

	update, err := tx.PrepareContext(ctx, updateQuery)
	if err != nil {
		return dataModels.Offering{}, err
	}
	defer update.Close()
	_, err = update.ExecContext(ctx, args...)
	if err != nil {

		return dataModels.Offering{}, err
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Offering{}, err
	}
	return ds.readOfferingByID(ctx, req.Row)

}

func (ds *OfferingDBDS) ListClassesTeacher(ctx context.Context, req offeringSchema.Pages, teacherID int64) (res []dataModels.TermClasses, total int, page int, err error) {
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
o.group_number AS OfferingGroupNumber,
c.course_number AS CourseNumber,
c.title AS title,
o.capacity AS capacity,
o.enrolled_count AS EnrolledCount,
o.week AS week,
o.day AS workDay, 
o.class_start_time AS ClassStartTime,
o.class_end_time AS ClassEndTime,
o.exam_start_time AS ExamStartTime,
o.exam_finish_time AS ExamFinishTime,
c.department_id AS department_id
FROM registration r
JOIN offerings o ON r.offering_row = o.row
JOIN courses c ON o.course_number = c.course_number
JOIN teachers t ON o.teacher_id = t.ID
WHERE o.teacher_id = ? AND o.term_id = ? AND o.isActive = 1
ORDER BY o.class_start_time LIMIT ? OFFSET ?;
`
	rows, err := ds.db.QueryContext(ctx, selectQuery, teacherID, termID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var class dataModels.DetailClasses
		err = rows.Scan(&class.OfferingGroupNumber, &class.CourseNumber, &class.Title, &class.Capacity, &class.EnrolledCount, &class.Week, &class.WorkDay, &class.ClassStartTime, &class.ClassEndTime, &class.ExamStartTime, &class.ExamFinishTime, &class.DepartmentID)
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

func (ds *OfferingDBDS) readOfferingByID(ctx context.Context, row int64) (res dataModels.Offering, err error) {
	var offering dataModels.Offering
	readQuery := fmt.Sprintf("SELECT row , group_number , course_number , teacher_id , capacity , enrolled_count , isActive , reserveation , term_id , week , day , class_start_time , class_end_time , exam_start_time , exam_finish_time , updated_at FROM %s WHERE row = ? ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, readQuery, row).Scan(&offering.Row, &offering.GroupNumber, &offering.CourseNumber, &offering.TeacherID, &offering.Capacity, &offering.EnrolledCount, &offering.IsActive, &offering.Reservation, &offering.TermID, &offering.Week, &offering.Day, &offering.ClassStartTime, &offering.ClassEndTime, &offering.ExamStartTime, &offering.ExamFinishTime, &offering.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return dataModels.Offering{}, errors.New(sql.ErrNoRows.Error())
		}
		return dataModels.Offering{}, err
	}
	return offering, nil
}

func (ds *OfferingDBDS) checkID(ctx context.Context, row int64) error {
	var check bool
	searchQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, searchQuery, row).Scan(&check)
	if err != nil {
		return errors.New(err.Error())
	}
	if !check {
		return errors.New("This is not a valid ID")
	}
	return nil
}
