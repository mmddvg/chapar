package http

import (
	"mmddvg/chapar/pkg/responses"
	"mmddvg/chapar/pkg/services/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *httpWs) GetChats(c echo.Context) error {
	res, err := h.App.GetChats(utils.GetUserId(c))
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, res)
}

func (h *httpWs) GetPvMessages(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(400, responses.Error{Message: "invalid id"})
	}

	res, err := h.App.GetPvMessages(id)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, res)
}

func (h *httpWs) GetGroupMessages(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(400, responses.Error{Message: "invalid id"})
	}

	res, err := h.App.GetGroupMessages(id)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, res)
}
