package main

import (
	HttpWs "mmddvg/chapar/pkg/http"
	"mmddvg/chapar/pkg/services"
)

func main() {
	app := services.NewApp()
	app.Run()
	HttpWs.Initiate(app)
}
