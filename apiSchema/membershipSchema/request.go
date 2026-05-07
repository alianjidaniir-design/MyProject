package membershipSchema

type CreateMembershipRequest struct {
	StudentID        int64  `json:"student_id"`
	ProgramRow       int64  `json:"program_row"`
	StatusMembership string `json:"status_membership"`
}
