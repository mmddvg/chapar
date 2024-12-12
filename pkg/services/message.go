package services

import "mmddvg/chapar/pkg/models"

func (app *Application) GetPvMessages(pvId uint64) ([]models.PvMessage, error) {
	return app.messageDB.GetPvMessages(pvId)
}

func (app *Application) GetGroupMessages(groupId uint64) ([]models.GroupMessage, error) {
	return app.messageDB.GetGroupMessages(groupId)
}
