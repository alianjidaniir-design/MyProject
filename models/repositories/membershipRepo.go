package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip"
	"context"
)

type MembershipRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.CreateMembershipRequest]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error)
}

var MemberShipRepo MembershipRepository = memberShip.GetRepo()
