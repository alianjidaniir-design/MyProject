package mySqlDS

import (
	"MyProject/apiSchema/departmentSchema"
	"MyProject/models/department/dataModels"
	"MyProject/models/department/dataSources"
	"MyProject/pkg/pagination"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DepartmentDBDS struct {
	tableSQL string
	db       *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		fmt.Println(err)
	}
	return loc
}

func NewDepartmentDBDS(tableName string, db *sql.DB) (dataSources.DepartmentDB, error) {
	ff := &DepartmentDBDS{
		tableSQL: tableName,
		db:       db,
	}
	return ff, nil
}
func (ds *DepartmentDBDS) CreateDepartment(ctx context.Context, req departmentSchema.CreateDepartmentReq) (dataModels.Department, error) {
	tx, err := ds.db.Begin()
	if err != nil {
		return dataModels.Department{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	var lastID int64
	err = tx.QueryRowContext(ctx, fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", ds.tableSQL)).Scan(&lastID)
	if err != nil {
		return dataModels.Department{}, errors.New("there was a problem coalesce")
	}
	newID := lastID + 1
	now := time.Now().In(myLocation())
	insertQuery := fmt.Sprintf(
		"INSERT INTO %s (id, college, educational_group, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		ds.tableSQL,
	)
	_, err = tx.ExecContext(ctx, insertQuery, newID, req.College, req.EducationalGroup, now, now)
	if err != nil {
		tx.Rollback()
		return dataModels.Department{}, fmt.Errorf("insert failed: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Department{}, errors.New("there was a problem trying to commit")
	}
	return ds.readDepartmentByID(ctx, newID)
}
func (ds *DepartmentDBDS) UpdateDepartment(ctx context.Context, req departmentSchema.UpdateDepartmentReq) (dataModels.Department, error) {
	now := time.Now().In(myLocation())
	err := ds.chackDepartment(ctx, req.ID)
	if err != nil {
		return dataModels.Department{}, errors.New("there is not department here")
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE id = ?", ds.tableSQL)
	update, err := ds.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return dataModels.Department{}, err
	}
	defer update.Close()
	result, err := update.ExecContext(ctx, now, req.ID)
	if err != nil {
		return dataModels.Department{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return dataModels.Department{}, err
	}

	return ds.readDepartmentByID(ctx, req.ID)
}

func (ds *DepartmentDBDS) ListDepartment(ctx context.Context, req departmentSchema.ListReq) ([]dataModels.Department, int64, error) {
	var departments []dataModels.Department
	page, pageSize, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return []dataModels.Department{}, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var total int64
	totalItem := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableSQL)
	err = ds.db.QueryRowContext(ctx, totalItem).Scan(&total)
	if err != nil {
		return []dataModels.Department{}, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT id , college , educational_group ,  created_at, updated_at FROM %s LIMIT ? OFFSET ?", ds.tableSQL)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []dataModels.Department{}, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var department dataModels.Department
		var createdAt, updatedAt sql.NullTime
		err = rows.Scan(&department.ID, &department.College, &department.EducationalGroup, &createdAt, &updatedAt)
		if err != nil {
			return []dataModels.Department{}, 0, err
		}
		if createdAt.Valid {
			department.CreatedAt = createdAt.Time.In(myLocation())
		}
		if updatedAt.Valid {
			department.UpdatedAt = updatedAt.Time.In(myLocation())
		}

		departments = append(departments, department)

	}

	if rows.Err() != nil {
		return []dataModels.Department{}, 0, err

	}

	return departments, total, nil

}

func (ds *DepartmentDBDS) readDepartmentByID(ctx context.Context, id int64) (dataModels.Department, error) {
	var department dataModels.Department
	readQuery := fmt.Sprintf("SELECT id ,college,educational_group , created_at , updated_at FROM %s WHERE id = ?", ds.tableSQL)
	var createdAt, updatedAt sql.NullTime
	if err := ds.db.QueryRowContext(ctx, readQuery, id).Scan(&department.ID, &department.College, &department.EducationalGroup, &createdAt, &updatedAt); err != nil {
		return dataModels.Department{}, err
	}
	if createdAt.Valid {
		department.CreatedAt = createdAt.Time
	} else {
		department.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		department.UpdatedAt = updatedAt.Time
	} else {
		department.UpdatedAt = time.Time{}
	}

	return department, nil
}

func (ds *DepartmentDBDS) DeleteCourse(ctx context.Context, req departmentSchema.DeleteDepartmentReq) (dataModels.Department, error) {
	var department dataModels.Department
	err := ds.chackDepartment(ctx, req.ID)
	if err != nil {
		return dataModels.Department{}, errors.New("Department Found not")
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=?", ds.tableSQL)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return dataModels.Department{}, err
	}
	return department, nil
}

func (ds *DepartmentDBDS) chackDepartment(ctx context.Context, ID int64) error {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM departments WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, search, ID).Scan(&check)

	if err != nil {
		return err
	}
	if !check {
		return errors.New("Department not found")
	}
	return nil
}
