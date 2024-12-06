package responses

import "mmddvg/chapar/pkg/models"

type Login struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

func NewLogin(user models.User, token string) Login {
	return Login{
		Id:       user.Id,
		Name:     user.Name,
		UserName: user.UserName,
		Token:    token,
	}
}
