package models

type UserDB interface {
	SignUp(NewUser) User
	Get(uint64) User
	GetContacts(uint64) []uint64
	IsContact(userId uint64, contactId uint64) bool
	AddContact(userId uint64, contactId uint64)
	RemoveContact(userId uint64, contactId uint64)
	Block(userId uint64, targetId uint64)
	UnBlock(userId uint64, targetId uint64)
	AddProfile(uint64, string) []string
	RemoveProfile(userId uint64, count uint) []string
	CreateGroup(uint64, title string, link string) Group
	AddGroupMember(userId uint64, memberId uint64)
	RemoveGroupMember(userId uint64, memberId uint64)
}
