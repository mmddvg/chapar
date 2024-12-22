package http

import (
	"mmddvg/chapar/pkg/requests"
	"mmddvg/chapar/pkg/services/utils"

	"github.com/labstack/echo/v4"
)

func (h *httpWs) SignUp(c echo.Context) error {

	body, err := GetBody[requests.User](c)
	if err != nil {
		return err
	}

	res, err := h.App.SignUp(body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(201, res)
}

func (h *httpWs) Login(c echo.Context) error {
	body, err := GetBody[requests.Login](c)
	if err != nil {
		return err
	}

	res, err := h.App.Login(body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, res)
}

func (h *httpWs) AddContact(c echo.Context) error {

	res, err := h.App.AddContact(utils.GetUserId(c), c.Param("contact_username"))
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(201, res)
}
func (h *httpWs) RemoveContact(c echo.Context) error {

	res, err := h.App.RemoveContact(utils.GetUserId(c), c.Param("contact_username"))
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, res)
}

func (h *httpWs) GetUser(c echo.Context) error {
	user, err := h.App.GetUser(utils.GetUserId(c))
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, user)
}

func (h *httpWs) GetContacts(c echo.Context) error {
	contacts, err := h.App.GetContacts(utils.GetUserId(c))
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, contacts)
}
