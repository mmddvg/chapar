package models

type UserDB interface {
	SignUp(NewUser) (User, error)
	Get(uint64) (User, error)
	GetContacts(uint64) ([]uint64, error)
	IsContact(userId uint64, contactId uint64) (bool, error)
	AddContact(userId uint64, contactId uint64) ([]uint64, error)
	RemoveContact(userId uint64, contactId uint64) ([]uint64, error)
	CreatePv(userId uint64, targetId uint64) (PrivateChat, error)
	GetPv(targetId uint64) (PrivateChat, error)
	Block(userId uint64, targetId uint64) (uint64, error)
	UnBlock(userId uint64, targetId uint64) (uint64, error)
	AddProfile(uint64, string) ([]string, error)
	RemoveProfile(userId uint64, count uint) ([]string, error)
	CreateGroup(uint64, title string, link string) (Group, error)
	AddGroupMember(userId uint64, memberId uint64) (GroupMember, error)
	RemoveGroupMember(userId uint64, memberId uint64) (GroupMember, error)
}

type MessageDB interface {
	WritePv(NewPvMessage) (PvMessage, error)
	WriteGroup(NewGroupMessage) (GroupMessage, error)
	EditPv(EditPvMessage) (PvMessage, error)
	EditGroup(EditGroupMessage) (GroupMessage, error)
}
