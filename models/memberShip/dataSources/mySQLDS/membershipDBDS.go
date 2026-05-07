package mySQLDS

import (
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip/dataModel"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type MembershipDBDS struct {
	tableName string
	db        *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*60*60+30*60)
	}
	return loc
}

func NewMembershipDBDS(tableName string, db *sql.DB) (*MembershipDBDS, error) {
	m := &MembershipDBDS{
		tableName: tableName,
		db:        db,
	}
	return m, nil
}

func (ds *MembershipDBDS) CreateMembership(ctx context.Context, req membershipSchema.CreateMembershipRequest) (res dataModel.Membership, err error) {
	var checkStudent bool
	checkStudentQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM student WHERE ID = ?)THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, checkStudentQuery, req.StudentID).Scan(&checkStudent)
	if err != nil {
		return dataModel.Membership{}, err
	}
	if !checkStudent {
		return dataModel.Membership{}, errors.New("student does not exist")
	}
	var checkProgram bool
	checkProgramQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM programs WHERE row = ? )THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, checkProgramQuery, req.ProgramRow).Scan(&checkProgram)
	if err != nil {
		return dataModel.Membership{}, err
	}
	if !checkProgram {
		return dataModel.Membership{}, errors.New("program does not exist")
	}
	var lastTime *time.Time
	now := time.Now().In(myLocation())
	if req.StatusMembership == constants.Approved {
		lastTime = &now
	} else if req.StatusMembership == constants.Rejected || req.StatusMembership == constants.Review {
		lastTime = nil
	} else {
		return res, errors.New("invalid status membership")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (student_id , program_row ,created_membership_at , status_membership , created_at , updated_at ) VALUES (?, ?, ? , ? ,? , ?)", ds.tableName)
	lastID, err := ds.db.ExecContext(ctx, insertQuery, req.StudentID, req.ProgramRow, lastTime, req.StatusMembership, now, now)
	if err != nil {
		return dataModel.Membership{}, err
	}
	insertID, err := lastID.LastInsertId()
	if err != nil {
		return dataModel.Membership{}, err
	}
	return ds.selectMembership(ctx, insertID)

}

func (ds *MembershipDBDS) selectMembership(ctx context.Context, ID int64) (res dataModel.Membership, err error) {
	var membership dataModel.Membership
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE ID = ?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&membership.ID, &membership.StudentID, &membership.ProgramRow, &membership.CreatedMemberShipAt, &membership.FinishMemberShipAt, &membership.StatusMemberShip, &membership.CreatedAt, &membership.UpdatedAt, &membership.DeletedAt)
	if err != nil {
		return membership, err
	}

	return membership, nil
}
