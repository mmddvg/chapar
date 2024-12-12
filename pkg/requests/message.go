package requests

import "mmddvg/chapar/pkg/models"

type Message struct {
	Reciever_id uint64 `json:"reciever_id" validate:"required"`
	ActionType  uint8  `json:"action_type" validate:"required,oneof= 0 1 2"`
	TargetType  uint8  `json:"target_type" validate:"required,oneof= 0 1"`
	Message     string `json:"message" validate:"required_if=ActionType 0 1"`
	MessageId   uint64 `json:"message_id" validate:"required_if=ActionType 1 2"`
}

func (m Message) RecieverId() uint64 {
	return m.Reciever_id
}

func (m Message) Action() models.ActionType {
	return models.ActionType(m.ActionType)
}

func (m Message) Target() models.TargetType {
	return models.TargetType(m.TargetType)
}
