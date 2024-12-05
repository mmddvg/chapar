package services

type ActionType uint8

const (
	NewMessage  ActionType = iota
	EditMessage ActionType = iota
)

type TargetType uint8

const (
	Pv    TargetType = iota
	Group TargetType = iota
)

type Message interface {
	RecieverId() uint64
	Action() ActionType
	Target() TargetType
}
