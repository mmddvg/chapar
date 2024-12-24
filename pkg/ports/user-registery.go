package ports

type UserRegister interface {
	Register(userId uint64, serverId uint) error
	Retrive(userId uint64) ([]uint, error)
	UnRegister(userId uint64, serverId uint) error
}
