package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher"
	"context"
)

type PublisherRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[publisherSchema.CreatePublisher]) (res publisherSchema.DetailPublisher, errStr string, code int, err error)
	Detail(ctx context.Context, req commonSchema.BaseRequest[publisherSchema.GetPublisher]) (res publisherSchema.DetailPublisher, errStr string, code int, err error)
}

var PublisherRepo PublisherRepository = publisher.GetRepo()
