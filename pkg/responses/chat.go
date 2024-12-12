package responses

import "mmddvg/chapar/pkg/models"

type ChatList struct {
	Users  []User         `json:"users"`
	Groups []models.Group `json:"groups"`
}

func NewChatList(users []models.User, groups []models.Group) ChatList {
	tmp := []User{}

	for _, v := range users {
		tmp = append(tmp, User{Id: v.Id, Name: v.Name, Username: v.UserName})
	}
	return ChatList{
		Users:  tmp,
		Groups: groups,
	}
}

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
