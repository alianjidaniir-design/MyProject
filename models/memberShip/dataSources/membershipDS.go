package dataSources

import (
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip/dataModel"
	"context"
)

type MembershipDS interface {
	CreateMembership(ctx context.Context, req membershipSchema.CreateMembershipRequest) (res dataModel.Membership, err error)
}
