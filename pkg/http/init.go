package http

import (
	"mmddvg/chapar/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpWs struct {
	App *services.Application
}

func Initiate(app *services.Application) {
	h := httpWs{App: app}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", h.Hello)
	e.GET("/message", h.chat)

	e.Logger.Fatal(e.Start(":8080"))
}
