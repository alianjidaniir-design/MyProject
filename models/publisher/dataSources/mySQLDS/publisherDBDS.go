package mySQLDS

import (
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher/dataModel"
	"MyProject/models/publisher/dataSources"
	"MyProject/pkg/pagination"
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
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	return ds.selected(ctx, req.ID)
}

func (ds *PublisherDBDS) DeletePublisher(ctx context.Context, req publisherSchema.GetPublisher) (dataModel.Publisher, error) {
	err := Val.CheckValidation(req)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return dataModel.Publisher{}, err
	}
	return dataModel.Publisher{}, nil
}

func (ds *PublisherDBDS) ListPublisher(ctx context.Context, req publisherSchema.PaginationPublisher) (res []dataModel.Publisher, total int, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return nil, 0, err
	}
	var publishers []dataModel.Publisher
	page, pageSize, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var tot int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&tot)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var p dataModel.Publisher
		err = rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Address)
		if err != nil {
			return nil, 0, err
		}
		publishers = append(publishers, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return publishers, tot, nil
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
