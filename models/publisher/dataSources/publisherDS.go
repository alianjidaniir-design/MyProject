package dataSources

import (
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher/dataModel"
	"context"
)

type PublisherDS interface {
	CreatePublisher(ctx context.Context, req publisherSchema.CreatePublisher) (res dataModel.Publisher, err error)
	DetailPublisher(ctx context.Context, req publisherSchema.GetPublisher) (res dataModel.Publisher, err error)
}
