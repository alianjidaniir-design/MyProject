package dataSources

import (
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip/dataModel"
	"context"
)

type MembershipDS interface {
	CreateMembership(ctx context.Context, req membershipSchema.CreateMembershipRequest) (res dataModel.Membership, err error)
	DeleteMembership(ctx context.Context, req membershipSchema.GetIDMembership) (res dataModel.Membership, err error)
	UpdateMembership(ctx context.Context, req membershipSchema.UpdateMembership) (res dataModel.Membership, err error)
	DeActiveMembership(ctx context.Context, req membershipSchema.GetIDMembership) (res dataModel.Membership, err error)
	DetailMembership(ctx context.Context, req membershipSchema.GetIDMembership) (res dataModel.Membership, err error)
}
