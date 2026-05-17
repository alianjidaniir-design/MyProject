package mySQLDS

import (
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher/dataModel"
	"MyProject/models/publisher/dataSources"
	Val "MyProject/pkg/val"
	"context"
	"errors"
	"fmt"

	"database/sql"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type PublisherDBDS struct {
	tableName string
	db        *sql.DB
}

func NewPublisherDBDS(tableName string, db *sql.DB) (dataSources.PublisherDS, error) {
	ff := &PublisherDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *PublisherDBDS) CreatePublisher(ctx context.Context, req publisherSchema.CreatePublisher) (res dataModel.Publisher, err error) {
	err = validate.Struct(req)
	if err != nil {
		var str string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, v := range validationErrors {
				str = Val.SwitchValidateErr(v)
				break
			}
			return dataModel.Publisher{}, fmt.Errorf("%v", str)

		}
		str = err.Error()
		return dataModel.Publisher{}, errors.New(str)

	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (name , phone , address) VALUES (?, ?, ?)", ds.tableName)
	result, err := ds.db.ExecContext(ctx, insertQuery, req.Name, req.Phone, req.Address)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return dataModel.Publisher{}, err
	}
	return ds.selected(ctx, lastInsertID)

}

func (ds *PublisherDBDS) selected(ctx context.Context, ID int64) (dataModel.Publisher, error) {
	var publisher dataModel.Publisher
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&publisher.Name, &publisher.Phone, &publisher.Address)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	return publisher, nil
}
