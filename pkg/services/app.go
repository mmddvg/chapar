package services

import "mmddvg/chapar/pkg/models"

type Application struct {
	// hub
	users      map[uint64]*Client
	channel    chan Message
	RegChan    chan Register
	UnregChann chan UnRegister

	// rest

	userDB    models.UserDB
	messageDB models.MessageDB
}

func NewApp(userDB models.UserDB, messageDB models.MessageDB) *Application {
	app := &Application{
		users:      make(map[uint64]*Client, 1000),
		channel:    make(chan Message, 1000),
		RegChan:    make(chan Register),
		UnregChann: make(chan UnRegister),

		userDB:    userDB,
		messageDB: messageDB,
	}

	app.Run()
	return app
}

func (app *Application) Hello() string {
	return "hello world"
}
