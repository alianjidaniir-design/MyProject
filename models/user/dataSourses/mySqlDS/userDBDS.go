package mySqlDS

import (
	"MyProject/apiSchema/userSchema"
	userDataModel "MyProject/models/user/dataModel"
	"context"
	"fmt"
)

type UserDBDS struct {
	tablename string
	tableSQL  string
	db        DBExecture
}

func (ds *UserDBDS) CreateStudent(ctx context.Context, req userSchema.LoginRequest) (userDataModel.User, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (code , name , family ) VALUES (?, ? , ?)", ds.tableSQL)
	insertResult, err := ds.db.ExecContext(ctx, insertQuery, req.Code, req.Name, req.Family)
	if err != nil {
		return userDataModel.User{}, err
	}
	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return userDataModel.User{}, err
	}
	return ds.readTaskByID(ctx, insertedID)
}

func (ds *UserDBDS) ReadStudent(ctx context.Context) (userDataModel.User, error) {
	var student userDataModel.User
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE ", ds.tableSQL)
	selectResult, err := ds.db.QueryContext(ctx, selectQuery)
	if err != nil {
		return userDataModel.User{}, err
	}
	return student, selectResult.Scan()
}

func (ds *UserDBDS) readTaskByID(ctx context.Context, userID int64) (userDataModel.User, error) {
	var students userDataModel.User
	readQuery := fmt.Sprintf("SELECT id , code , name , family FROM %s WHERE id = ?", ds.tableSQL)
	if err := ds.db.QueryRowContext(ctx, readQuery, userID).Scan(&students.ID, &students.Code, students.Name, students.Family); err != nil {
		return userDataModel.User{}, err
	}
	return students, nil

}

func (ds *UserDBDS) TableName() string {
	return ds.tablename
}
