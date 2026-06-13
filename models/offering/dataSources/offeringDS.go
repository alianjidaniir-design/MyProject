package dataSources

import (
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering/dataModels"
	"context"
)

type OfferingDS interface {
	CreateOffering(ctx context.Context, req offeringSchema.CreateOfferingRequest) (res dataModels.Offering, err error)
	ListOffering(ctx context.Context, req offeringSchema.ListOfferingsRequest) (res []dataModels.ListOfferings, total int, err error)
	GetOffering(ctx context.Context, req offeringSchema.GetRowOfferingRequest) (res dataModels.Offering, err error)
	DeActiveOffering(ctx context.Context, req offeringSchema.GetRowOfferingRequest) (res dataModels.Offering, err error)
	EditOffering(ctx context.Context, req offeringSchema.EditOffering) (res dataModels.Offering, err error)
}
