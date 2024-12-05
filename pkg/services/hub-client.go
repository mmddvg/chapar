package services

import (
	"github.com/google/uuid"
)

type Client struct {
	devices map[uuid.UUID]chan Message
}

func NewClient(id uuid.UUID, ch chan Message) *Client {
	return &Client{
		devices: map[uuid.UUID]chan Message{
			id: ch,
		},
	}
}

type Register struct {
	Id    uint64
	UId   uuid.UUID
	Write chan Message
}

type UnRegister struct {
	Id  uint64
	UId uuid.UUID
}
