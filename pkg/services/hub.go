package services

import (
	"log/slog"
	"mmddvg/chapar/pkg/models"
	"os"
	"strconv"
)

func (h *Application) Run() {
	slog.Info("starting the hub")
	go func() {
		for {
			select {
			case re := <-h.RegChan:
				h.register(re)
			case unre := <-h.UnregChann:
				h.unregister(unre)
			case m := <-h.channel:
				h.handleMessage(m)
			}
		}
	}()
	go h.registerUserGlobally()
	go h.unregisterUserGlobally()
}
func (h *Application) handleMessage(message models.HubMessage) {
	switch message.TargetType {
	case models.GroupTarget:
		h.handleGroupMessage(message)
	case models.PvTarget:
		h.handleSingleMessage(message)
	}
}

func (h *Application) handleSingleMessage(message models.HubMessage) {
	if v, ok := h.users[message.RecieverId]; ok {
		for id := range v.devices {
			v.devices[id] <- message // todo : make sure this doesn't get blocked
		}
	}

	if v, ok := h.users[message.SenderId]; ok {
		for id := range v.devices {
			v.devices[id] <- message
		}
	}
}

// todo : handle in separate goroutine ?
func (h *Application) handleGroupMessage(message models.HubMessage) {

	ids, err := h.userDB.GetGroupMembers(message.RecieverId)
	if err != nil {
		_ = err.Error()
		return
	}

	for _, id := range ids {
		if v, ok := h.users[id]; ok {
			for _, ch := range v.devices {
				ch <- message
			}
		}
	}
}

func (h *Application) register(arg Register) {
	slog.Info("register locally ", "arg", arg)
	v, ok := h.users[arg.Id]
	if !ok {
		h.globalRegister <- arg.Id
		h.users[arg.Id] = NewClient(arg.UId, arg.Write)
	} else {
		v.devices[arg.UId] = arg.Write
	}
}

func (h *Application) unregister(arg UnRegister) {
	v, ok := h.users[arg.Id]
	if ok {
		delete(v.devices, arg.UId)
	}
	if len(v.devices) == 0 {
		delete(h.users, arg.Id)
		h.globalUnRegister <- arg.Id
	}
}

// when a user wants to write a message , the message is sent to `Channel` channel to be broadcasted
// when a user should read a message (which is writing from prespective of server app) it listens to channels on `devices` map

func (h *Application) registerUserGlobally() {
	tmp, err := strconv.ParseUint(os.Getenv("SERVER_ID"), 10, 0)
	if err != nil {
		slog.Error("couldn't parse server id ", err)
	}

	serverId := uint(tmp)
	for userId := range h.globalRegister {
		slog.Info("globally registering user", "id", userId)
		err = h.userRegister.Register(userId, serverId)
		if err != nil {
			slog.Error("error globally registering user  ", "err", err.Error())
		}
	}
}

func (h *Application) unregisterUserGlobally() {
	tmp, err := strconv.ParseUint(os.Getenv("SERVER_ID"), 10, 0)
	if err != nil {
		slog.Error("couldn't parse server id ", err)
	}

	serverId := uint(tmp)
	for userId := range h.globalUnRegister {
		slog.Info("globally unregistering user", "id", userId)

		err = h.userRegister.UnRegister(userId, serverId)
		if err != nil {
			slog.Error("error globally unregistering user  ", "err", err.Error())
		}
	}
}
