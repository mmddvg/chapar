package responses

import (
	"mmddvg/chapar/pkg/models"

	"github.com/samber/lo"
)

type ChatList struct {
	Pvs    []PvInfo       `json:"pvs"`
	Groups []models.Group `json:"groups"`
}

func NewChatList(userId uint64, users []models.PrivateChat, groups []models.Group) ChatList {
	tmp := []PvInfo{}

	for _, v := range users {
		tmp = append(tmp, PvInfo{Id: v.Id, UserId: lo.Ternary(v.User1 == userId, v.User2, v.User1)})
	}
	return ChatList{
		Pvs:    tmp,
		Groups: groups,
	}
}

type PvInfo struct {
	Id     uint64 `json:"id"`
	UserId uint64 `json:"user_id"`
}
