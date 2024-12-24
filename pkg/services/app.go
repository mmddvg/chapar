package services

import (
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/ports"
)

type Application struct {
	// hub
	users            map[uint64]*Client
	channel          chan models.HubMessage
	RegChan          chan Register
	UnregChann       chan UnRegister
	globalRegister   chan uint64 // these two for handling global user registration like in redis
	globalUnRegister chan uint64

	// rest

	userDB    ports.UserDB
	messageDB ports.MessageDB

	profileStorage ports.PictureStorage
	userRegister   ports.UserRegister
}

func NewApp(userDB ports.UserDB, messageDB ports.MessageDB, profileStorage ports.PictureStorage, userRegister ports.UserRegister) *Application {
	app := &Application{
		users:      make(map[uint64]*Client, 1000),
		channel:    make(chan models.HubMessage, 1000),
		RegChan:    make(chan Register),
		UnregChann: make(chan UnRegister),

		userDB:           userDB,
		messageDB:        messageDB,
		profileStorage:   profileStorage,
		userRegister:     userRegister,
		globalRegister:   make(chan uint64, 100),
		globalUnRegister: make(chan uint64, 100),
	}

	app.Run()
	return app
}

func (app *Application) Hello() string {
	return "hello world"
}
