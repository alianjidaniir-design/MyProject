package mySQLDS

import (
	"MyProject/apiSchema/subjectSchema"
	"MyProject/models/subject/dataModel"
	"MyProject/models/subject/dataSources"
	"MyProject/pkg/pagination"
	Val "MyProject/pkg/val"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SubjectDBDS struct {
	tableName string
	db        *sql.DB
}

func NewSubjectDBDS(tableName string, db *sql.DB) (dataSources.SubjectDS, error) {
	ff := &SubjectDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *SubjectDBDS) CreateSubject(ctx context.Context, req subjectSchema.CreateSubject) (res dataModel.Subject, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Subject{}, err
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (name , description) VALUES (?, ?)", ds.tableName)
	result, err := ds.db.ExecContext(ctx, insertQuery, req.Name, req.Description)
	if err != nil {
		return dataModel.Subject{}, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return dataModel.Subject{}, err
	}
	return ds.selectDetailSubject(ctx, lastId)
}

func (ds *SubjectDBDS) GetSubject(ctx context.Context, req subjectSchema.GetSubject) (res dataModel.Subject, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Subject{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Subject{}, err
	}
	return ds.selectDetailSubject(ctx, req.ID)
}

func (ds *SubjectDBDS) ListSubjects(ctx context.Context, req subjectSchema.Pagination) (res []dataModel.Subject,total int, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return res, 0, err
	}
	var subjects []dataModel.Subject
	page , pageSize , err:=pagination.CheckPage(req.Page , req.Size)
	if err != nil {
		return nil ,0, err
	}
	limit:=pageSize
	offset:=(page - 1)*limit
	var tot int
	countQuery:=fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&tot)
	if err != nil {
		return nil ,0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery , limit, offset)
	if err != nil {
		return nil ,0, err
	}
	defer rows.Close()
	for rows.Next() {
		var subject dataModel.Subject
		err=rows.Scan(&subject.ID, &subject.Name, &subject.Description)
		if err != nil {
			return nil ,0, err
		}
		subjects = append(subjects, subject)
	}
	err = rows.Err()
	if err != nil {
		return nil ,0, err
	}
	return subjects , tot, err
}


func (ds *SubjectDBDS) selectDetailSubject(ctx context.Context, ID int64) (dataModel.Subject, error) {
	var subject dataModel.Subject
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&subject.ID, &subject.Name, &subject.Description)
	if err != nil {
		return dataModel.Subject{}, err
	}
	return subject, nil
}

func (ds *SubjectDBDS) DeleteSubject(ctx context.Context, req subjectSchema.GetSubject) (res dataModel.Subject, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Subject{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Subject{}, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return dataModel.Subject{}, err
	}
	return dataModel.Subject{}, nil
}

func (ds *SubjectDBDS) checkID(ctx context.Context, id int64) error {
	var check bool
	checkQuery := `
SELECT EXISTS (SELECT 1 FROM subjects WHERE id=?)`
	err := ds.db.QueryRowContext(ctx, checkQuery, id).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("This is ID there is not")
	}
	return nil
}
