package dataSources

import (
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering/dataModels"
	"context"
)

type OfferingDS interface {
	CreateOffering(ctx context.Context , req  offeringSchema.CreateOfferingRequest)(res dataModels.Offering , err error )
}