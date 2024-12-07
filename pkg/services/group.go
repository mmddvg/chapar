package services

import (
	"mime/multipart"
	"mmddvg/chapar/pkg/errs"
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/requests"
)

func (app *Application) CreateGroup(ownerId uint64, body requests.NewGroup) (models.Group, error) {
	return app.userDB.CreateGroup(ownerId, body.Name, body.Link)
}

func (app *Application) AddGroupMember(userId uint64, body requests.Member) (models.GroupMember, error) {
	var (
		err error
		res models.GroupMember
	)
	isCon, err := app.userDB.IsContact(body.MemberId, userId)
	if err != nil {
		return res, err
	}
	if !isCon {
		return res, errs.NewBadRequest("not a contact")
	}

	return app.userDB.AddGroupMember(body.GroupId, body.MemberId)
}

func (app *Application) RemoveGroupMember(userId uint64, body requests.Member) (models.GroupMember, error) {
	var (
		err error
		res models.GroupMember
	)

	group, err := app.userDB.GetGroup(body.GroupId)
	if err != nil {
		return res, err
	}

	if group.OwnerId != userId {
		return res, errs.NewBadRequest("only owner can remove")
	}

	return app.userDB.RemoveGroupMember(body.GroupId, body.MemberId)
}

func (app *Application) UpdateGroup(userId uint64, body requests.UpdateGroup) (models.Group, error) {
	var (
		err error
		res models.Group
	)

	group, err := app.userDB.GetGroup(body.GroupId)
	if err != nil {
		return res, err
	}

	if group.OwnerId != userId {
		return res, errs.NewBadRequest("only owner can update")
	}

	return app.userDB.UpdateGroup(body)
}

func (app *Application) AddGroupProfile(userId uint64, groupId uint64, file multipart.File, contentType string) (models.GroupProfile, error) {
	var (
		err error
		res models.GroupProfile
	)
	group, err := app.userDB.GetGroup(groupId)
	if err != nil {
		return res, err
	}

	if group.OwnerId != userId {
		return res, errs.NewBadRequest("only owner can add profile")
	}

	link, err := app.profileStorage.Save(file, contentType)
	if err != nil {
		return res, err
	}

	res, err = app.userDB.AddGroupProfile(groupId, link)
	if err != nil {
		tmp := app.profileStorage.Delete(link)
		if tmp != nil {
			return res, tmp
		} else {
			return res, err
		}
	}

	return res, err
}

func (app *Application) RmGroupProfile(userId uint64, body requests.RmGroupProfile) error {
	group, err := app.userDB.GetGroup(body.GroupId)
	if err != nil {
		return err
	}

	if group.OwnerId != userId {
		return errs.NewBadRequest("only owner can add profile")
	}

	uid, err := app.userDB.RmGroupProfile(body)
	if err != nil {
		return err
	}

	return app.profileStorage.Delete(uid)

}
