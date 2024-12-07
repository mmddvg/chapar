package requests

type NewGroup struct {
	Name string `json:"name" validate:"required"`
	Link string `json:"link" validate:"required"`
}

type Member struct {
	MemberId uint64 `json:"member_id" validate:"required"`
	GroupId  uint64 `json:"group_id" validate:"required"`
}

type UpdateGroup struct {
	GroupId uint64 `json:"group_id" validte:"required"`
	Name    string `json:"name" validte:"required"`
}

type RmGroupProfile struct {
	GroupId  uint64 `json:"group_id" validate:"required"`
	NthCount uint8  `json:"nth_count" validate:"required"`
}
