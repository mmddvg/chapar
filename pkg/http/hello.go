package http

import "github.com/labstack/echo/v4"

func (e *httpWs) Hello(c echo.Context) error {
	return c.String(200, e.App.Hello())
}
