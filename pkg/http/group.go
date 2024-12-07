package http

import (
	"mmddvg/chapar/pkg/requests"
	"mmddvg/chapar/pkg/responses"
	"mmddvg/chapar/pkg/services/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *httpWs) CreateGroup(c echo.Context) error {
	body, err := GetBody[requests.NewGroup](c)
	if err != nil {
		return err
	}

	res, err := h.App.CreateGroup(utils.GetUserId(c), body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(201, res)
}

func (h *httpWs) AddGroupMember(c echo.Context) error {
	body, err := GetBody[requests.Member](c)
	if err != nil {
		return err
	}

	res, err := h.App.AddGroupMember(utils.GetUserId(c), body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(201, res)
}

func (h *httpWs) RmGroupMember(c echo.Context) error {
	body, err := GetBody[requests.Member](c)
	if err != nil {
		return err
	}

	res, err := h.App.RemoveGroupMember(utils.GetUserId(c), body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(200, res)
}

func (h *httpWs) AddGroupProfile(c echo.Context) error {
	groupId, err := strconv.ParseUint(c.Param("group_id"), 10, 0)
	if err != nil {
		return c.JSON(400, responses.Error{Message: "invalid group id "})
	}
	src, contentType, err := GetFile(c)
	if err != nil {
		return err
	}
	defer src.Close()

	res, err := h.App.AddGroupProfile(utils.GetUserId(c), groupId, src, contentType)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(201, res)
}

func (h *httpWs) RmGroupProfile(c echo.Context) error {
	body, err := GetBody[requests.RmGroupProfile](c)
	if err != nil {
		return err
	}
	err = h.App.RmGroupProfile(utils.GetUserId(c), body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.String(204, "")
}

func (h *httpWs) UpdateGroup(c echo.Context) error {
	body, err := GetBody[requests.UpdateGroup](c)
	if err != nil {
		return err
	}

	res, err := h.App.UpdateGroup(utils.GetUserId(c), body)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(201, res)
}
