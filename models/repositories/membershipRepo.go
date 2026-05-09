package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip"
	"context"
)

type MembershipRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.CreateMembershipRequest]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.GetIDMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.UpdateMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error)
	DeActive(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.GetIDMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.GetIDMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.PaginationMemberShip]) (res membershipSchema.ListMembershipSchema, errStr string, code int, err error)
}

var MemberShipRepo MembershipRepository = memberShip.GetRepo()
