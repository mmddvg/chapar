package models

type User struct {
	Id       uint64 `db:"id"`
	Name     string `db:"name"`
	UserName string `db:"username"`
	Password string `db:"string"`
	Contacts []uint64
	Profiles []string
	Blocked  []uint64
}

type NewUser struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
