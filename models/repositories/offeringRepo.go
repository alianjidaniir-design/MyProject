package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering"
	"context"
)

type OfferingRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.CreateOfferingRequest]) (res offeringSchema.CreateOfferingResponse, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.ListOfferingsRequest]) (res offeringSchema.ListOfferingResponse, errStr string, code int, err error)
}

var OfferingRepo OfferingRepository = offering.GetRepository()
