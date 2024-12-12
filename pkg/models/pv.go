package models

import "time"

type PrivateChat struct {
	Id        uint64    `db:"id" json:"id"`
	User1     uint64    `db:"user1" json:"user1"`
	User2     uint64    `db:"user2" json:"user2"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
