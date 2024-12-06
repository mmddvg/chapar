package http

import (
	"mmddvg/chapar/pkg/requests"

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
