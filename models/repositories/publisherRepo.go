package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/publisherSchema"
	"context"
)

type PublisherRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[publisherSchema.CreatePublisher]) (res publisherSchema.DetailPublisher, errStr string, code int, err error)
}

var PublisherRepo PublisherRepository
