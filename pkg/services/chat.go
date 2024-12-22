package services

import (
	"mmddvg/chapar/pkg/responses"
)

func (app *Application) GetChats(userId uint64) (responses.ChatList, error) {

	users, groups, err := app.userDB.GetChats(userId)
	if err != nil {
		return responses.ChatList{}, err
	}

	return responses.NewChatList(userId, users, groups), nil
}
