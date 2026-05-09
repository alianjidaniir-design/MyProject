package membershipSchema

type CreateMembershipRequest struct {
	StudentID        int64  `json:"student_id"`
	ProgramRow       int64  `json:"program_row"`
	StatusMembership string `json:"status_membership"`
}

type GetIDMembership struct {
	ID int64 `json:"id"`
}

type UpdateMembership struct {
	ID               int64  `json:"id"`
	StatusMembership string `json:"status_membership"`
}
