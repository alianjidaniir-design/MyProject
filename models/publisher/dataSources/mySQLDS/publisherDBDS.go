package mySQLDS

import (
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher/dataModel"
	"MyProject/models/publisher/dataSources"
	Val "MyProject/pkg/val"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

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
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Publisher{}, err
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

func (ds *PublisherDBDS) DetailPublisher(ctx context.Context, req publisherSchema.GetPublisher) (res dataModel.Publisher, err error) {
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	return ds.selected(ctx, req.ID)
}

func (ds *PublisherDBDS) selected(ctx context.Context, ID int64) (dataModel.Publisher, error) {
	var publisher dataModel.Publisher
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&publisher.ID, &publisher.Name, &publisher.Phone, &publisher.Address)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	return publisher, nil
}

func (ds *PublisherDBDS) checkID(ctx context.Context, ID int64) error {
	var check bool
	selectQuery := `
SELECT EXISTS (SELECT 1 FROM ` + ds.tableName + ` WHERE id=?)`
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("ID not found")
	}
	return nil
}
