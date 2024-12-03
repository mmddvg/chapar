package models

type UserDB interface {
	SignUp(NewUser) User
	Get(uint64) User
	GetContacts(uint64) []uint64
	AddContact(userId uint64, contactId uint64)
	RemoveContact(userId uint64, contactId uint64)
	Block(userId uint64, targetId uint64)
	UnBlock(userId uint64, targetId uint64)
	AddProfile(uint64, string) []string
	RemoveProfile(userId uint64, count uint) []string
}
