package services

import (
	"log/slog"
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/requests"
)

type Message interface {
	RecieverId() uint64
	Action() models.ActionType
	Target() models.TargetType
}

func (h *Application) SendMessage(userId uint64, m requests.Message) {
	if m.Target() == models.GroupTarget {
		isMember, err := h.userDB.IsMember(m.RecieverId(), userId)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if !isMember {
			slog.Info("user trying to send a message in a group they are not a member")
			return
		}

		if m.Action() == models.NewMessage {
			_, err = h.messageDB.WriteGroup(models.NewGroupMessage{
				GroupId:  m.RecieverId(),
				Message:  m.Message,
				SenderId: userId,
			})
		} else if m.Action() == models.EditMessage {
			_, err = h.messageDB.EditGroup(models.EditGroupMessage{
				Id:         m.MessageId,
				GroupId:    m.RecieverId(),
				NewMessage: m.Message,
			})
		}

		if err != nil {
			_ = err.Error()
			return
		}

	} else if m.Target() == models.PvTarget {
		// check if is blocked and create a pv if there isn't one
		isBlocked, err := h.userDB.IsBlocked(m.RecieverId(), userId)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if isBlocked {
			slog.Info("user trying to send messaged while blocked")
			return
		}

		pv, err := h.userDB.GetPvOrCreate(userId, m.RecieverId())
		if err != nil {
			_ = err.Error()
			return
		}

		if m.Action() == models.NewMessage {
			_, err = h.messageDB.WritePv(models.NewPvMessage{
				PvId:     pv.Id,
				SenderId: userId,
				Message:  m.Message,
			})
		} else if m.Action() == models.EditMessage {
			_, err = h.messageDB.EditPv(models.EditPvMessage{
				Id:         m.MessageId,
				NewMessage: m.Message,
			})
		} else if m.Action() == models.SeenAck {
			_, err = h.messageDB.SeenAck(m.MessageId)
		}

		if err != nil {
			_ = err.Error()
			return
		}
	}

	h.channel <- m
}
