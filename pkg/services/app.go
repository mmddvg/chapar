package services

import (
	"mmddvg/chapar/pkg/ports"
)

type Application struct {
	// hub
	users      map[uint64]*Client
	channel    chan Message
	RegChan    chan Register
	UnregChann chan UnRegister

	// rest

	userDB    ports.UserDB
	messageDB ports.MessageDB

	profileStorage ports.PictureStorage
}

func NewApp(userDB ports.UserDB, messageDB ports.MessageDB, profileStorage ports.PictureStorage) *Application {
	app := &Application{
		users:      make(map[uint64]*Client, 1000),
		channel:    make(chan Message, 1000),
		RegChan:    make(chan Register),
		UnregChann: make(chan UnRegister),

		userDB:         userDB,
		messageDB:      messageDB,
		profileStorage: profileStorage,
	}

	app.Run()
	return app
}

func (app *Application) Hello() string {
	return "hello world"
}
