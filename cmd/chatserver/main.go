package main

import (
	"log"
	HttpWs "mmddvg/chapar/pkg/http"
	"mmddvg/chapar/pkg/repositories/postgres"
	"mmddvg/chapar/pkg/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	db := postgres.NewPostgresRepo()
	app := services.NewApp(db, db)

	HttpWs.Initiate(app)
}
