package dataSources

import (
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher/dataModel"
	"context"
)

type PublisherDS interface {
	CreatePublisher(ctx context.Context, req publisherSchema.CreatePublisher) (res dataModel.Publisher, err error)
	DetailPublisher(ctx context.Context, req publisherSchema.GetPublisher) (res dataModel.Publisher, err error)
	DeletePublisher(ctx context.Context, req publisherSchema.GetPublisher) (dataModel.Publisher, error)
	ListPublisher(ctx context.Context, req publisherSchema.PaginationPublisher) (res []dataModel.Publisher, total int, err error)
}
