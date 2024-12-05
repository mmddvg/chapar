package services

type Application struct {
	// hub
	users      map[uint64]*Client
	channel    chan Message
	RegChan    chan Register
	UnregChann chan UnRegister

	// rest

}

func NewApp() *Application {
	return &Application{
		users:      make(map[uint64]*Client, 1000),
		channel:    make(chan Message, 1000),
		RegChan:    make(chan Register),
		UnregChann: make(chan UnRegister),
	}
}

func (app *Application) Hello() string {
	return "hello world"
}
