package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetBody[T any](c echo.Context) (T, error) {
	var req T
	if err := c.Bind(&req); err != nil {
		return req, echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = e.ActualTag()
		}
		return req, echo.NewHTTPError(http.StatusBadRequest, validationErrors)
	}

	return req, nil
}
