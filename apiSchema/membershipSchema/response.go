package membershipSchema

import "MyProject/models/memberShip/dataModel"

type DetailMembershipSchema struct {
	MemberShip dataModel.Membership `json:"member_ship"`
	Massage    string               `json:"massage"`
}

type ListMembershipSchema struct {
	List  []dataModel.Membership `json:"member_ship"`
	Total int                    `json:"total"`
}
