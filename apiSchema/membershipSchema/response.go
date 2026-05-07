package membershipSchema

import "MyProject/models/memberShip/dataModel"

type DetailMembershipSchema struct {
	MemberShip dataModel.Membership `json:"member_ship"`
	Massage    string               `json:"massage"`
}
