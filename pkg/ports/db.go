package ports

import (
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/requests"
)

type UserDB interface {
	SignUp(requests.User) (models.User, error)
	Get(uint64) (models.User, error)
	GetByUsername(string) (models.User, error)
	GetContacts(uint64) ([]uint64, error)
	IsContact(userId uint64, contactId uint64) (bool, error)
	AddContact(userId uint64, contactId uint64) ([]uint64, error)
	RemoveContact(userId uint64, contactId uint64) ([]uint64, error)
	CreatePv(userId uint64, targetId uint64) (models.PrivateChat, error)
	GetPvOrCreate(userId, targetId uint64) (models.PrivateChat, error)
	Block(userId uint64, targetId uint64) (uint64, error)
	IsBlocked(userId uint64, targetId uint64) (bool, error)
	UnBlock(userId uint64, targetId uint64) (uint64, error)
	AddProfile(uint64, string) ([]string, error)
	RemoveProfile(userId uint64, count uint) ([]string, error)
	CreateGroup(ownerId uint64, title string, link string) (models.Group, error)
	GetGroup(uint64) (models.Group, error)
	GetGroupMembers(uint64) ([]uint64, error)
	IsMember(groupId uint64, userId uint64) (bool, error)
	UpdateGroup(requests.UpdateGroup) (models.Group, error)
	AddGroupProfile(uint64, string) (models.GroupProfile, error)
	RmGroupProfile(requests.RmGroupProfile) (string, error)
	AddGroupMember(groupId uint64, memberId uint64) (models.GroupMember, error)
	RemoveGroupMember(groupId uint64, memberId uint64) (models.GroupMember, error)

	GetChats(uint64) ([]models.User, []models.Group, error)
}

type MessageDB interface {
	WritePv(models.NewPvMessage) (models.PvMessage, error)
	WriteGroup(models.NewGroupMessage) (models.GroupMessage, error)
	EditPv(models.EditPvMessage) (models.PvMessage, error)
	EditGroup(models.EditGroupMessage) (models.GroupMessage, error)
	SeenAck(uint64) (models.PvMessage, error)

	GetPvMessages(uint64) ([]models.PvMessage, error)
	GetGroupMessages(uint64) ([]models.GroupMessage, error)
}
