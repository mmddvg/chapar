package models

import "time"

type User struct {
	Id       uint64 `db:"id"`
	Name     string `db:"name"`
	UserName string `db:"username"`
}

type UserProfile struct {
	UserId    uint64    `db:"user_id"`
	Link      string    `db:"link"`
	CreatedAt time.Time `db:"created_at"`
}

type Contact struct {
	UserId    uint64 `db:"user_id"`
	ContactId uint64 `db:"contact_id"`
}

type Blocked struct {
	UserId   uint64 `db:"user_id"`
	TargetId uint64 `db:"target_id"`
}
