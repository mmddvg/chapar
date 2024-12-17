package models

import (
	"database/sql"
	"time"

	"github.com/samber/lo"
)

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
	Id        uint64       `db:"id" json:"id"`
	PvId      uint64       `db:"pv_id" json:"pv_id"`
	SenderId  uint64       `db:"sender_id" json:"sender_id"`
	Message   string       `db:"message" json:"message"`
	SeenAt    sql.NullTime `db:"seen_at" json:"seen_at"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
}

type NewPvMessage struct {
	PvId     uint64 `json:"pv_id"`
	SenderId uint64 `json:"sender_id"`
	Message  string `json:"message"`
}

type GroupMessage struct {
	Id        uint64    `db:"id" json:"id"`
	GroupId   uint64    `db:"group_id" json:"group_id"`
	Message   string    `db:"message" json:"message"`
	SenderId  uint64    `db:"sender_id" json:"sender_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
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

type HubMessage struct {
	Id         uint64     `json:"id"`
	RecieverId uint64     `json:"reciever_id"`
	ChatId     uint64     `json:"chat_id"`
	SenderId   uint64     `json:"sender_id"`
	Message    string     `json:"message"`
	SeenAt     *time.Time `json:"seen_at,omitempty"`
	TargetType TargetType `json:"target_type"`
	ActionType ActionType `json:"action_type"`
	CreatedAt  time.Time  `json:"created_at"`
}

// reciverId is only used for private messages
func NewHubMessage(message any, recieverId uint64, actionType ActionType) HubMessage {
	switch v := message.(type) {
	case PvMessage:
		return HubMessage{
			Id:         v.Id,
			RecieverId: recieverId,
			ChatId:     v.PvId,
			SenderId:   v.SenderId,
			Message:    v.Message,
			SeenAt:     lo.Ternary(v.SeenAt.Valid, &v.SeenAt.Time, nil),
			TargetType: PvTarget,
			ActionType: actionType,
			CreatedAt:  v.CreatedAt,
		}
	case GroupMessage:
		return HubMessage{
			Id:         v.Id,
			ChatId:     v.GroupId,
			RecieverId: v.GroupId,
			SenderId:   v.SenderId,
			Message:    v.Message,
			TargetType: GroupTarget,
			ActionType: actionType,
			CreatedAt:  v.CreatedAt,
		}
	}

	panic("unhandled message type") // shouldn't panic in a real world app
}
