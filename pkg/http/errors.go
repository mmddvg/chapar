package http

import (
	"errors"
	"log/slog"
	"mmddvg/chapar/pkg/errs"
	"mmddvg/chapar/pkg/responses"

	"github.com/labstack/echo/v4"
)

func ErrHandler(c echo.Context, err error) error {
	res := responses.Error{Message: err.Error()}
	if errors.As(err, &errs.ErrBadRequest{}) || errors.As(err, &errs.ErrDuplicate{}) {
		return c.JSON(400, res)
	} else if errors.As(err, &errs.ErrNotFound{}) {
		return c.JSON(404, res)
	} else if errors.As(err, &errs.ErrUnexpected{}) {
		return c.JSON(500, res)
	}

	slog.Error("uncaught error : " + err.Error())
	return c.JSON(500, responses.Error{Message: "unexpected error"})
}
