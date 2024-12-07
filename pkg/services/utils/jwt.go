package utils

import (
	"fmt"
	"os"
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

func ValidateJWT(token string) (jwt.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken.Claims, nil
}
