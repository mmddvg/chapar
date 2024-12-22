package models

type User struct {
	Id       uint64 `db:"id"`
	Name     string `db:"name"`
	UserName string `db:"username"`
	Password string `db:"password"`
	Contacts []uint64
	Profiles []string
	Blocked  []uint64
}

type Contact struct {
	ContactId       uint64  `db:"contact_id" json:"contact_id"`
	ContactUsername string  `db:"contact_username" json:"contact_username"`
	PvId            *uint64 `db:"pv_id" json:"pv_id"`
}
