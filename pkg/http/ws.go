package http

import (
	"log/slog"
	"mmddvg/chapar/pkg/errs"
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/requests"
	"mmddvg/chapar/pkg/responses"
	"mmddvg/chapar/pkg/services"
	"mmddvg/chapar/pkg/services/utils"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func (h *httpWs) chat(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	id, err := wsAuth(ws)
	if err != nil {
		ws.WriteJSON(responses.Error{Message: err.Error()})
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid message format"))
		return nil
	}

	uid, err := uuid.NewV7()
	if err != nil {
		slog.Error("error creating connection uid: ", err)
		return err
	}

	ch := make(chan models.HubMessage)
	h.App.RegChan <- services.Register{
		Id:    id,
		UId:   uid,
		Write: ch,
	}

	go func() {
		for {
			var inp requests.Message
			err := ws.ReadJSON(&inp)
			if err != nil {
				slog.Error("read error: ", err)
				break
			}
			h.App.SendMessage(id, inp)
		}
	}()

	for m := range ch {
		err := ws.WriteJSON(m)
		if err != nil {
			slog.Error("write error: ", err)
			break
		}
	}

	h.App.UnregChann <- services.UnRegister{Id: id, UId: uid}
	return nil
}

func wsAuth(ws *websocket.Conn) (uint64, error) {
	var err error
	var tokenMessage struct {
		Token string `json:"token"`
	}

	err = ws.ReadJSON(&tokenMessage)
	if err != nil {
		slog.Error("failed to read token message: ", err)
		return 0, errs.NewBadRequest("network error")
	}

	claims, err := utils.ValidateJWT(tokenMessage.Token)
	if err != nil {
		return 0, errs.NewBadRequest("invalid jwt")
	}

	// TODO
	// unreachable errors , but better be handled
	idStr, _ := claims.GetSubject()

	id, _ := strconv.ParseUint(idStr, 10, 0)

	return id, nil
}
