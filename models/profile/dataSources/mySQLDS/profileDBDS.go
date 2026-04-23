package mySQLDS

import (
	"MyProject/apiSchema/profileSchema"
	"MyProject/models/profile/dataModels"
	"MyProject/pkg/pagination"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ProfileDBDS struct {
	tableName string
	db        *sql.DB
}

func NewProfileDBDS(tableName string, db *sql.DB) (*ProfileDBDS, error) {
	cfg := &ProfileDBDS{
		tableName: tableName,
		db:        db,
	}
	return cfg, nil
}

func (ds *ProfileDBDS) CreateScoreStudent(ctx context.Context, req profileSchema.CreateScoresReq) (res dataModels.Profile, err error) {
	var profile dataModels.Profile
	tx, err := ds.db.BeginTx(ctx, nil)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	var checking bool
	checkRegisterQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE ID = ? AND deleted_at IS NULL AND status = ? ) THEN 1 ELSE 0 END
`
	var enroll = constants.StatusEnrolled
	err = tx.QueryRowContext(ctx, checkRegisterQuery, req.RegistrationID, enroll).Scan(&checking)
	if err != nil {
		return dataModels.Profile{}, err
	}
	if !checking {
		return dataModels.Profile{}, errors.New("Registration does not exist")
	}
	insertQuery := `
INSERT INTO profiles (registration_id , status_score , grade , score) VALUES (?, ?, ?, ?)
`
	switch {
	case req.Score > 20 || req.Score < 0:
		return dataModels.Profile{}, errors.New("Score must be between 0 and 20")
	case req.Score > 17:
		profile.Grade = "A"
		profile.StatusScore = constants.Passed
	case req.Score > 14:
		profile.Grade = "B"
		profile.StatusScore = constants.Passed
	case req.Score > 10:
		profile.Grade = "C"
		profile.StatusScore = constants.Passed
	case req.Score > 7:
		profile.Grade = "D"
		profile.StatusScore = constants.Failed
	case req.Score < 7:
		profile.Grade = "E"
		profile.StatusScore = constants.Failed
	default:
		return dataModels.Profile{}, errors.New("invalid score")
	}
	lastID, err := tx.ExecContext(ctx, insertQuery, req.RegistrationID, profile.StatusScore, profile.Grade, req.Score)
	if err != nil {
		return dataModels.Profile{}, err
	}

	result, err := lastID.LastInsertId()
	if err != nil {
		return dataModels.Profile{}, err
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Profile{}, err
	}
	return ds.readOProfileByID(ctx, result)
}

func (ds *ProfileDBDS) ListScoresStudents(crx context.Context, req profileSchema.ListAllScoresReq) (res []dataModels.ScoresStudents, total int, err error) {
	var profile []dataModels.ScoresStudents
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalRows int
	countQuery := `SELECT COUNT(*) FROM profiles`
	err = ds.db.QueryRowContext(crx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := `
SELECT
u.ID    AS student_id,
u.code AS student_code,
c.ID AS course_id,
c.course_number AS course_number,
o.row AS offering_row,
o.group_number AS offering_group_number,
o.teacher_id AS offering_teacher_id,
p.status_score ,
p.grade ,
p.score 
FROM profiles p 
JOIN registration r ON p.registration_id = r.ID 
JOIN offerings o ON r.offering_row = o.row 
JOIN courses c ON o.course_id = c.ID 
JOIN student u ON r.student_id = u.ID 
ORDER BY u.code LIMIT ? OFFSET ?;
`
	rows, err := ds.db.QueryContext(crx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var s dataModels.ScoresStudents
		err = rows.Scan(
			&s.StudentID,
			&s.StudentCode,
			&s.CourseID,
			&s.CourseNumber,
			&s.OfferingRows,
			&s.OfferingGroup,
			&s.OfferingTeacher,
			&s.StatusScore,
			&s.Grade,
			&s.Score,
		)
		if err != nil {
			return nil, 0, err
		}
		profile = append(profile, s)

	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return profile, totalRows, nil
}

func (ds *ProfileDBDS) readOProfileByID(ctx context.Context, ID int64) (res dataModels.Profile, err error) {
	var profile dataModels.Profile
	readQuery := fmt.Sprintf("SELECT ID , registration_id , status_score , grade  , score FROM %s WHERE ID = ? ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&profile.ID, &profile.RegistrationID, &profile.StatusScore, &profile.Grade, &profile.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			return dataModels.Profile{}, errors.New(sql.ErrNoRows.Error())
		}
		return dataModels.Profile{}, err
	}
	return profile, nil
}

func (ds *ProfileDBDS) checkID(ctx context.Context, row int64) error {
	var check bool
	searchQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM profiles WHERE row = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, searchQuery, row).Scan(&check)
	if err != nil {
		return errors.New(err.Error())
	}
	if !check {
		return errors.New("This is not a valid Row")
	}
	return nil
}
