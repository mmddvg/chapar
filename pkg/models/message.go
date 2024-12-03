package models

import "time"

type PvMessage struct {
	Id        uint64    `db:"id"`
	PvId      uint64    `db:"pv_id"`
	SenderId  uint64    `db:"sender_id"`
	Message   string    `db:"message"`
	SeenAt    time.Time `db:"seen_at"`
	CreatedAt time.Time `db:"created_at"`
}
