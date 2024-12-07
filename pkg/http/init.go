package http

import (
	"mmddvg/chapar/pkg/services"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpWs struct {
	App *services.Application
}

func Initiate(app *services.Application) {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	}

	h := httpWs{App: app}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/signup", h.SignUp)
	e.POST("/login", h.Login)

	r := e.Group("/restricted", echojwt.WithConfig(config))

	r.GET("/hello", h.Hello)
	r.GET("/message", h.chat)

	g := r.Group("/group")

	g.POST("", h.CreateGroup)
	g.PATCH("", h.UpdateGroup)
	g.PUT("/member", h.AddGroupMember)
	g.DELETE("/member", h.RmGroupMember)
	g.PUT("/profile/:group_id", h.AddGroupProfile)
	g.DELETE("/profile", h.RmGroupProfile)

	e.Logger.Fatal(e.Start(":8080"))
}
