package mySqlDS

import (
	"MyProject/apiSchema/termSchema"
	termDataModel "MyProject/models/term/dataModels"
	"MyProject/models/term/dataSources"
	"MyProject/pkg/pagination"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TermDBDS struct {
	tableSQL string
	db       *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func NewTermDBDS(tableName string, db *sql.DB) (dataSources.TermDS, error) {
	ff := &TermDBDS{
		tableSQL: tableName,
		db:       db,
	}
	return ff, nil
}

func (ds *TermDBDS) CreateTerm(ctx context.Context, req termSchema.CreateTerm) (res termDataModel.Term, err error) {
	if req.Term < 1 || req.Term > 8 {
		return termDataModel.Term{}, errors.New("term is invalid . that should be between 1 and 8")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (term , year) VALUES (?,?)", ds.tableSQL)
	result, err := ds.db.Exec(insertQuery, req.Term, req.Year)
	if err != nil {
		return termDataModel.Term{}, errors.New(err.Error())
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return termDataModel.Term{}, errors.New("error getting last id from database")
	}
	return ds.readTermByID(ctx, lastId)
}

func (ds *TermDBDS) DeleteTerms(ctx context.Context, req termSchema.DeleteTerm) (res termDataModel.Term, err error) {
	err = ds.chackUser(ctx, req.ID)
	if err != nil {
		return termDataModel.Term{}, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=? ", ds.tableSQL)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return termDataModel.Term{}, errors.New(err.Error())
	}
	return termDataModel.Term{}, nil
}

func (ds *TermDBDS) ListTerms(ctx context.Context, req termSchema.ListTerm) (res []termDataModel.Term, total int, err error) {
	var terms []termDataModel.Term
	page, pageSize, err := pagination.CheckPage(req.PageIndex, req.PageSize)
	if err != nil {
		return []termDataModel.Term{}, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalAll int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableSQL)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&totalAll)
	if err != nil {
		return []termDataModel.Term{}, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT ID, term , year FROM %s ORDER BY ID LIMIT ? OFFSET ?", ds.tableSQL)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []termDataModel.Term{}, 0, errors.New("pagination error")
	}
	defer rows.Close()
	for rows.Next() {
		var term termDataModel.Term
		err = rows.Scan(&term.ID, &term.Term, &term.Year)
		if err != nil {
			return []termDataModel.Term{}, 0, errors.New("Scaning with error")
		}
		terms = append(terms, term)
	}
	if err = rows.Err(); err != nil {
		return []termDataModel.Term{}, 0, err
	}
	return terms, totalAll, err

}

func (ds *TermDBDS) readTermByID(ctx context.Context, termID int64) (termDataModel.Term, error) {
	var term termDataModel.Term
	readQuery := fmt.Sprintf("SELECT ID , term , year FROM %s WHERE ID = ?", ds.tableSQL)

	if err := ds.db.QueryRowContext(ctx, readQuery, termID).Scan(&term.ID, &term.Term, &term.Year); err != nil {
		return termDataModel.Term{}, err
	}

	return term, nil

}
func (ds *TermDBDS) chackUser(ctx context.Context, ID int64) error {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM terms WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, search, ID).Scan(&check)

	if err != nil {
		return err
	}
	if !check {
		return errors.New("Term not found")
	}
	return nil
}
