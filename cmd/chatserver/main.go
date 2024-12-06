package main

import (
	HttpWs "mmddvg/chapar/pkg/http"
	"mmddvg/chapar/pkg/repositories/postgres"
	"mmddvg/chapar/pkg/services"
)

func main() {
	db := postgres.NewPostgresRepo()
	app := services.NewApp(db, db)

	HttpWs.Initiate(app)
}
