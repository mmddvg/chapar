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
