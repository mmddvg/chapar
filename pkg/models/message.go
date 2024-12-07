package models

import "time"

type ActionType uint8

const (
	NewMessage  ActionType = iota
	EditMessage ActionType = iota
	SeenAck     ActionType = iota
)

type TargetType uint8

const (
	PvTarget    TargetType = iota
	GroupTarget TargetType = iota
)

type PvMessage struct {
	Id        uint64    `db:"id"`
	PvId      uint64    `db:"pv_id"`
	SenderId  uint64    `db:"sender_id"`
	Message   string    `db:"message"`
	SeenAt    time.Time `db:"seen_at"`
	CreatedAt time.Time `db:"created_at"`
}

type NewPvMessage struct {
	PvId     uint64 `json:"pv_id"`
	SenderId uint64 `json:"sender_id"`
	Message  string `json:"message"`
}

type GroupMessage struct {
	Id        uint64 `db:"id"`
	GroupId   uint64 `db:"group_id"`
	Message   string `db:"message"`
	SenderId  uint64 `db:"sender_id"`
	CreatedAt uint64 `db:"created_at"`
}

type NewGroupMessage struct {
	GroupId  uint64 `db:"group_id"`
	Message  string `db:"message"`
	SenderId uint64 `db:"sender_id"`
}

type EditPvMessage struct {
	Id         uint64 `json:"id"`
	NewMessage string `json:"new_message"`
}

type EditGroupMessage struct {
	Id         uint64 `json:"id"`
	GroupId    uint64 `json:"group_id"`
	NewMessage string `json:"new_message"`
}
