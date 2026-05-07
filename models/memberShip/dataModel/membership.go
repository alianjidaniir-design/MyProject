package dataModel

import "time"

type Membership struct {
	ID                  int64      `json:"id"`
	StudentID           int64      `json:"student_id"`
	ProgramRow          int64      `json:"program_row"`
	CreatedMemberShipAt *time.Time `json:"created_membership_at"`
	FinishMemberShipAt  *time.Time `json:"finish_membership_at"`
	StatusMemberShip    string     `json:"status_membership"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}
