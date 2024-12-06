package http

import "github.com/labstack/echo/v4"

func ErrHandler(c echo.Context, err error) error {
	return c.JSON(400, err.Error())
}
