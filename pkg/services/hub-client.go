package services

import (
	"mmddvg/chapar/pkg/models"

	"github.com/google/uuid"
)

type Client struct {
	devices map[uuid.UUID]chan models.HubMessage
}

func NewClient(id uuid.UUID, ch chan models.HubMessage) *Client {
	return &Client{
		devices: map[uuid.UUID]chan models.HubMessage{
			id: ch,
		},
	}
}

type Register struct {
	Id    uint64
	UId   uuid.UUID
	Write chan models.HubMessage
}

type UnRegister struct {
	Id  uint64
	UId uuid.UUID
}
