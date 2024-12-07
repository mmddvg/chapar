package utils

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetUserId(c echo.Context) uint64 {
	t := c.Get("user").(*jwt.Token)
	s, _ := t.Claims.GetSubject()
	id, _ := strconv.ParseUint(s, 10, 0)
	return id
}
