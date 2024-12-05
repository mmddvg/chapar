package services

func (h *Application) Run() {
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
}
func (h *Application) handleMessage(message Message) {
	switch message.Target() {
	case Group:
		h.handleGroupMessage(message)
	case Pv:
		h.handleSingleMessage(message)
	}
}

func (h *Application) handleSingleMessage(message Message) {
	if v, ok := h.users[message.RecieverId()]; ok {
		for id := range v.devices {
			v.devices[id] <- message // todo : make sure this doesn't get blocked
		}
	}
}

func (h *Application) handleGroupMessage(message Message) {
	// todo

	ids := []uint64{}

	for _, id := range ids {
		if v, ok := h.users[id]; ok {
			for _, ch := range v.devices {
				ch <- message
			}
		}
	}
}

func (h *Application) register(arg Register) {
	v, ok := h.users[arg.Id]
	if !ok {
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
	}
}

func (h *Application) SendMessage(m Message) {
	h.channel <- m
}

// when a user wants to write a message , the message is sent to `Channel` channel to be broadcasted
// when a user should read a message (which is writing from prespective of server app) it listens to channels on `devices` map
