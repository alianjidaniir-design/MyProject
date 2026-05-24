package mySQLDS

import (
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip/dataModel"
	"MyProject/pkg/pagination"
	TimeLoc "MyProject/pkg/timeLoc"
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
	now := time.Now().In(TimeLoc.MyLocation())
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

func (ds *MembershipDBDS) DeleteMembership(ctx context.Context, req membershipSchema.GetIDMembership) (dataModel.Membership, error) {
	err := ds.checkID(ctx, req.ID)
	now := time.Now().In(TimeLoc.MyLocation())
	if err != nil {
		return dataModel.Membership{}, err
	}
	membership, err := ds.selectMembership(ctx, req.ID)
	if membership.DeletedAt != nil {
		return dataModel.Membership{}, errors.New("membership deleted before")
	}
	deleted := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE ID = ?", ds.tableName)
	row, err := ds.db.PrepareContext(ctx, deleted)
	if err != nil {
		return dataModel.Membership{}, err
	}
	defer row.Close()
	_, err = row.ExecContext(ctx, now, req.ID)
	if err != nil {
		return dataModel.Membership{}, err
	}

	return ds.selectMembership(ctx, req.ID)
}

func (ds *MembershipDBDS) UpdateMembership(ctx context.Context, req membershipSchema.UpdateMembership) (dataModel.Membership, error) {
	err := ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Membership{}, err
	}
	now := time.Now().In(TimeLoc.MyLocation())
	var createMembershipAt *time.Time
	var status string
	detailMembership, err := ds.selectMembership(ctx, req.ID)
	if detailMembership.DeletedAt != nil {
		return dataModel.Membership{}, errors.New("membership deleted before")
	}
	if req.StatusMembership == constants.Approved {
		if detailMembership.StatusMemberShip == constants.Approved {
			return dataModel.Membership{}, errors.New("membership status is already approved")
		}
		status = constants.Approved
		createMembershipAt = &now

	} else if req.StatusMembership == constants.Review {
		if detailMembership.StatusMemberShip == constants.Review {
			return dataModel.Membership{}, errors.New("membership status is already reviewed")
		} else if detailMembership.StatusMemberShip == constants.Approved {
			return dataModel.Membership{}, errors.New("can not update approved status to review status")
		}
		status = constants.Review

	} else if req.StatusMembership == constants.Rejected {
		if detailMembership.StatusMemberShip == constants.Rejected {
			return dataModel.Membership{}, errors.New("membership status is already rejected")
		} else if detailMembership.StatusMemberShip == constants.Approved {
			return dataModel.Membership{}, errors.New("can not update approved status to rejected status")
		}
		status = constants.Rejected

	} else {
		return dataModel.Membership{}, errors.New("invalid status membership")
	}
	update := fmt.Sprintf("UPDATE %s SET status_membership = ? , created_membership_at = ? , updated_at = ?  WHERE ID = ? ", ds.tableName)
	row, err := ds.db.PrepareContext(ctx, update)
	if err != nil {
		return dataModel.Membership{}, err
	}
	defer row.Close()
	_, err = row.ExecContext(ctx, status, createMembershipAt, now, req.ID)
	if err != nil {
		return dataModel.Membership{}, err
	}
	return ds.selectMembership(ctx, req.ID)
}
func (ds *MembershipDBDS) DeActiveMembership(ctx context.Context, req membershipSchema.GetIDMembership) (res dataModel.Membership, err error) {
	var membership dataModel.Membership
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return membership, err
	}
	membership, err = ds.selectMembership(ctx, req.ID)
	if membership.DeletedAt != nil {
		return membership, errors.New("membership deleted before")
	} else if membership.StatusMemberShip != constants.Approved {
		return membership, errors.New("membership is not approved")
	} else if membership.FinishMemberShipAt != nil {
		return membership, errors.New("membership finished before")
	}
	now := time.Now().In(TimeLoc.MyLocation())
	update := fmt.Sprintf("UPDATE `%s` SET updated_at = ? , finish_membership_at = ? WHERE ID = ?", ds.tableName)
	row, err := ds.db.PrepareContext(ctx, update)
	if err != nil {
		return membership, err
	}
	defer row.Close()
	_, err = row.ExecContext(ctx, now, now, req.ID)
	if err != nil {
		return membership, err
	}
	return ds.selectMembership(ctx, req.ID)
}

func (ds *MembershipDBDS) DetailMembership(ctx context.Context, req membershipSchema.GetIDMembership) (res dataModel.Membership, err error) {
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Membership{}, err
	}
	return ds.selectMembership(ctx, req.ID)
}

func (ds *MembershipDBDS) ListMembership(ctx context.Context, req membershipSchema.PaginationMemberShip) (res []dataModel.Membership, total int, err error) {
	var memberships []dataModel.Membership
	page, pageSize, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var membership dataModel.Membership
		err = rows.Scan(&membership.ID, &membership.StudentID, &membership.ProgramRow, &membership.CreatedMemberShipAt, &membership.FinishMemberShipAt, &membership.StatusMemberShip, &membership.CreatedAt, &membership.UpdatedAt, &membership.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		memberships = append(memberships, membership)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return memberships, totalRows, nil
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

func (ds *MembershipDBDS) checkID(ctx context.Context, ID int64) error {
	var check bool
	checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM memberships WHERE ID = ?)THEN 1 ELSE 0 END`
	err := ds.db.QueryRowContext(ctx, checkQuery, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("membership does not exist")
	}
	return nil
}
