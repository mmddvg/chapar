package main

import (
	"log"
	HttpWs "mmddvg/chapar/pkg/http"
	"mmddvg/chapar/pkg/repositories/localfs"
	"mmddvg/chapar/pkg/repositories/postgres"
	"mmddvg/chapar/pkg/repositories/redisregistery"
	"mmddvg/chapar/pkg/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	userRegister, err := redisregistery.NewRedisRegister()
	if err != nil {
		log.Fatal(err)
	}
	profileStore := localfs.NewLocalPictureStorage("./static")

	db := postgres.NewPostgresRepo()
	app := services.NewApp(db, db, profileStore, userRegister)

	HttpWs.Initiate(app)
}
