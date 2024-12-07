package models

import "time"

type Group struct {
	Id      uint64 `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	Link    string `db:"link" json:"link"`
	OwnerId uint64 `db:"owner_id" json:"owner_id"`
}

type GroupProfile struct {
	GId       uint64    `db:"g_id" json:"g_id"`
	Link      string    `db:"link" json:"link"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type GroupMember struct {
	GroupId  uint64    `db:"group_id" json:"group_id"`
	MemberId uint64    `db:"member_id" json:"member_id"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
}
