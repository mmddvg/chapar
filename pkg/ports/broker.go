package ports

type Broker interface {
	Send()
	Get()
}
