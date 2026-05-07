package mySQLDS

import (
	"MyProject/apiSchema/programSchema"
	"MyProject/models/program/dataModel"
	"MyProject/models/program/dataSources"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ProgramDBDS struct {
	tableName string
	db        *sql.DB
}

func NewProgramDBDS(tableName string, db *sql.DB) (dataSources.ProgramDS, error) {
	ff := &ProgramDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *ProgramDBDS) CreateProgram(ctx context.Context, req programSchema.CreateProgramRequest) (res dataModel.Program, err error) {
	var lastInsertID int64
	var check bool
	reviewCategory := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM categories WHERE row = ? )THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, reviewCategory, req.CategoryID).Scan(&check)
	if err != nil {
		return dataModel.Program{}, err
	}
	if !check {
		return dataModel.Program{}, errors.New("category does not exist")
	}
	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row) , 0 ) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, lastIDQuery).Scan(&lastInsertID)
	if err != nil {
		return dataModel.Program{}, err
	}
	newID := lastInsertID + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , category_row , name , description ) VALUES (?, ? , ? , ?)", ds.tableName)
	_, err = ds.db.ExecContext(ctx, insertQuery, newID, req.CategoryID, req.Name, req.Description)
	if err != nil {
		return dataModel.Program{}, err
	}
	return ds.selectProgram(ctx, newID)

}

func (ds *ProgramDBDS) GetProgram(ctx context.Context, req programSchema.GetDetailProgramRequest) (res dataModel.Program, err error) {
	err = ds.checkID(ctx, req.Row)
	if err != nil {
		return dataModel.Program{}, err
	}
	return ds.selectProgram(ctx, req.Row)
}

func (ds *ProgramDBDS) selectProgram(ctx context.Context, ID int64) (dataModel.Program, error) {
	var program dataModel.Program
	selectQuery := fmt.Sprintf("SELECT row , category_row , name , description FROM %s WHERE row = ?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&program.Row, &program.CategoryID, &program.Name, &program.Description)
	if err != nil {
		return dataModel.Program{}, err
	}
	return program, nil
}

func (ds *ProgramDBDS) checkID(ctx context.Context, ID int64) error {
	var check bool
	checkQuery := `
SELECT EXISTS (SELECT 1 FROM programs WHERE row = ? )`
	err := ds.db.QueryRowContext(ctx, checkQuery, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("program does not exist")
	}
	return nil
}
