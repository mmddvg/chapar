package services

type Application struct {
	// hub
	users      map[uint64]*Client
	Channel    chan Message
	RegChan    chan Register
	UnregChann chan UnRegister

	// rest

}

func NewApp() *Application {
	return &Application{
		users:      make(map[uint64]*Client, 1000),
		Channel:    make(chan Message, 1000),
		RegChan:    make(chan Register),
		UnregChann: make(chan UnRegister),
	}
}
