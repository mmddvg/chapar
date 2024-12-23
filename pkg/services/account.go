package services

import (
	"mmddvg/chapar/pkg/errs"
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/requests"
	"mmddvg/chapar/pkg/responses"
	"mmddvg/chapar/pkg/services/utils"
)

func (app *Application) SignUp(body requests.User) (responses.Login, error) {
	var (
		err error
		tmp responses.Login
	)

	body.Password, err = utils.Encrypt(body.UserName, body.Password)
	if err != nil {
		return tmp, err
	}

	user, err := app.userDB.SignUp(body)
	if err != nil {
		return tmp, err
	}

	token, err := utils.GenerateJWT(user.Id)
	if err != nil {
		return tmp, err
	}

	return responses.NewLogin(user, token), nil
}

func (app *Application) Login(body requests.Login) (responses.Login, error) {
	var (
		err error
		tmp responses.Login
	)
	user, err := app.userDB.GetByUsername(body.UserName)
	if err != nil {
		return tmp, err
	}

	isMatch, err := utils.CheckPassword(user.Password, body.Password, user.UserName)
	if err != nil {
		return tmp, err
	}

	if !isMatch {
		return tmp, errs.NewBadRequest("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.Id)
	if err != nil {
		return tmp, err
	}

	return responses.NewLogin(user, token), nil
}

func (app *Application) AddContact(userId uint64, contactUsername string) ([]models.Contact, error) {
	user, err := app.userDB.GetByUsername(contactUsername)
	if err != nil {
		return []models.Contact{}, err
	}

	return app.userDB.AddContact(userId, user.Id)
}

func (app *Application) RemoveContact(userId uint64, contactUsername string) ([]models.Contact, error) {
	user, err := app.userDB.GetByUsername(contactUsername)
	if err != nil {
		return []models.Contact{}, err
	}

	return app.userDB.RemoveContact(userId, user.Id)
}

func (app *Application) GetUser(userId uint64) (responses.Login, error) {
	user, err := app.userDB.Get(userId)
	if err != nil {
		return responses.Login{}, err
	}

	return responses.NewLogin(user, ""), nil
}

func (app *Application) GetContacts(userId uint64) ([]models.Contact, error) {
	return app.userDB.GetContacts(userId)
}
