package http

import (
	"log/slog"
	"mmddvg/chapar/pkg/services"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func (h *httpWs) chat(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	id, _ := strconv.ParseUint(c.QueryParam("id"), 10, 0)

	uid, err := uuid.NewV7()

	if err != nil {
		slog.Error("error creating connection uid : ", err)
	}
	ch := make(chan services.Message)
	h.App.RegChan <- services.Register{
		Id:    id,
		UId:   uid,
		Write: ch,
	}

	go func() {
		for {
			var inp WsMessage
			err := ws.ReadJSON(&inp)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			h.App.SendMessage(inp)
		}
	}()

	for m := range ch {
		err := ws.WriteJSON(m)
		if err != nil {
			slog.Error(err.Error())
		}

	}

	return c.String(200, "")
}

type WsMessage struct {
	Id      uint64 `json:"id"`
	Message string `json:"message"`
}

func (m WsMessage) Action() services.ActionType {
	return services.NewMessage
}

func (m WsMessage) Target() services.TargetType {
	return services.Pv
}

func (m WsMessage) RecieverId() uint64 {
	return m.Id
}
