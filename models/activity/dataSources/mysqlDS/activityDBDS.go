package mysqlDS

import (
	"MyProject/apiSchema/activitySchema"
	"MyProject/models/activity/dataModel"
	adminDataModels "MyProject/models/admins/dataModels"
	TimeLoc "MyProject/pkg/timeLoc"
	Val "MyProject/pkg/val"
	"fmt"

	"MyProject/models/activity/dataSources"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ActivityDBDS struct {
	tableName string
	db        *sql.DB
}

func NewActivityDBDS(tableName string, db *sql.DB) (dataSources.ActivityDS, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	ff := &ActivityDBDS{
		tableName: tableName,
		db:        db,
	}

	return ff, nil

}

func (ds *ActivityDBDS) CreateActivity(ctx context.Context, req activitySchema.CreateActivity, role string, ID int64) (res dataModel.Activity, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Activity{}, err
	}
	var checkTeacher, checkCourse bool
	ExistTeacher := `
SELECT EXISTS (SELECT 1 FROM teachers WHERE ID = ?)
SELECT EXISTS (SELECT 1 FROM courses WHERE ID = ?)
`
	err = ds.db.QueryRow(ExistTeacher, ID).Scan(&checkTeacher)
	if err != nil {
		return dataModel.Activity{}, err
	}
	if !checkTeacher {
		return adminDataModels.Admins{}, errors.New("there is not TeacherID")
	}

	now := time.Now().In(TimeLoc.MyLocation())

	insert := fmt.Sprintf("INSERT INTO %s (user_name , password , name , family , email , role_name , created_at ) VALUES (?, ?, ?, ?, ?, ?, ?)", ds.tableName)
	insertQuery, err := ds.db.ExecContext(ctx, insert, req.Username, hashing, req.Name, req.Family, req.Email, req.RoleName, now)
	if err != nil {
		return adminDataModels.Admins{}, err
	}

	insertID, err := insertQuery.LastInsertId()
	if err != nil {
		return adminDataModels.Admins{}, err
	}
	return ds.readQuery(ctx, insertID)

}

func (ds *ActivityDBDS) readQuery(ctx context.Context, ID int64) (adminDataModels.Admins, error) {
	var admin adminDataModels.Admins
	read := fmt.Sprintf("SELECT ID , user_name , password , name, family , email , role_name ,  created_at FROM %s WHERE ID=? ", ds.tableName)

	err := ds.db.QueryRowContext(ctx, read, ID).Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Name, &admin.Family, &admin.Email, &admin.RoleName, &admin.CreatedAt)
	if err != nil {
		return admin, err
	}

	return admin, nil
}
