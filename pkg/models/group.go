package models

import "time"

type Group struct {
	Id      uint64 `db:"id"`
	Title   string `db:"title"`
	Link    string `db:"link"`
	OwnerId uint64 `db:"owner_id"`
}

type GroupProfile struct {
	GId       uint64    `db:"g_id"`
	Link      string    `db:"link"`
	CreatedAt time.Time `db:"created_at"`
}

type GroupMember struct {
	GroupId   uint64    `db:"group_id"`
	MemberId  uint64    `db:"member_id"`
	JoinedAt  time.Time `db:"joined_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
