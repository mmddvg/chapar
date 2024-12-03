package models

type User struct {
	Id       uint64 `db:"id"`
	Name     string `db:"name"`
	UserName string `db:"username"`
	Contacts []uint64
	Profiles []string
	Blocked  []uint64
}

type NewUser struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
}
