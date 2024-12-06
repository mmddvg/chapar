package postgres

import (
	"database/sql"
	"errors"
	"mmddvg/chapar/pkg/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrInternal       = errors.New("internal error")
)

func (r *Repository) SignUp(newUser models.NewUser) (models.User, error) {
	query := `INSERT INTO users (username, name, password) VALUES ($1, $2, $3) RETURNING id, username, name`
	var user models.User
	err := r.db.Get(&user, query, newUser.UserName, newUser.Name, newUser.Password)
	if err != nil {
		if isDuplicateError(err) {
			return models.User{}, ErrDuplicateEntry
		}
		return models.User{}, ErrInternal
	}
	return user, nil
}

func (r *Repository) Get(userID uint64) (models.User, error) {
	query := `SELECT id, username, name FROM users WHERE id = $1`
	var user models.User
	err := r.db.Get(&user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, ErrInternal
	}
	return user, nil
}

func (r *Repository) GetContacts(userID uint64) ([]uint64, error) {
	query := `SELECT contact_id FROM contacts WHERE user_id = $1`
	var contacts []uint64
	err := r.db.Select(&contacts, query, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return contacts, nil
}

func (r *Repository) IsContact(userID, contactID uint64) (bool, error) {
	query := `SELECT COUNT(*) > 0 FROM contacts WHERE user_id = $1 AND contact_id = $2`
	var exists bool
	err := r.db.Get(&exists, query, userID, contactID)
	if err != nil {
		return false, ErrInternal
	}
	return exists, nil
}

func (r *Repository) AddContact(userID, contactID uint64) ([]uint64, error) {
	query := `INSERT INTO contacts (user_id, contact_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, userID, contactID)
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicateEntry
		}
		return nil, ErrInternal
	}
	return r.GetContacts(userID)
}

func (r *Repository) RemoveContact(userID, contactID uint64) ([]uint64, error) {
	query := `DELETE FROM contacts WHERE user_id = $1 AND contact_id = $2`
	_, err := r.db.Exec(query, userID, contactID)
	if err != nil {
		return nil, ErrInternal
	}
	return r.GetContacts(userID)
}

func (r *Repository) CreatePv(userID, targetID uint64) (models.PrivateChat, error) {
	query := `
		INSERT INTO private_chats (user1, user2) 
		VALUES ($1, $2) 
		RETURNING id, user1, user2, created_at`
	var chat models.PrivateChat
	err := r.db.Get(&chat, query, userID, targetID)
	if err != nil {
		if isDuplicateError(err) {
			return models.PrivateChat{}, ErrDuplicateEntry
		}
		return models.PrivateChat{}, ErrInternal
	}
	return chat, nil
}

func (r *Repository) GetPv(targetID uint64) (models.PrivateChat, error) {
	query := `SELECT id, user1, user2, created_at FROM private_chats WHERE id = $1`
	var chat models.PrivateChat
	err := r.db.Get(&chat, query, targetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PrivateChat{}, ErrNotFound
		}
		return models.PrivateChat{}, ErrInternal
	}
	return chat, nil
}

func (r *Repository) WritePv(newMessage models.NewPvMessage) (models.PvMessage, error) {
	query := `INSERT INTO pv_messages VALUES ($1, $2, $3, $4, $5) RETURNING *;`
	var message models.PvMessage
	err := r.db.Get(&message, query, generateId(), newMessage.PvId, newMessage.SenderId, newMessage.Message, time.Now())
	if err != nil {
		if isDuplicateError(err) {
			return models.PvMessage{}, ErrDuplicateEntry
		}
		return models.PvMessage{}, ErrInternal
	}
	return message, nil
}

func (r *Repository) WriteGroup(newMessage models.NewGroupMessage) (models.GroupMessage, error) {
	query := `
		INSERT INTO group_messages 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, group_id, message, sender_id, created_at`
	var message models.GroupMessage
	err := r.db.Get(&message, query, generateId(), newMessage.GroupId, newMessage.Message, newMessage.SenderId, time.Now())
	if err != nil {
		if isDuplicateError(err) {
			return models.GroupMessage{}, ErrDuplicateEntry
		}
		return models.GroupMessage{}, ErrInternal
	}
	return message, nil
}

func (r *Repository) EditPv(edit models.EditPvMessage) (models.PvMessage, error) {
	query := `UPDATE pv_messages SET message = $1 WHERE id = $2 RETURNING *;`
	var message models.PvMessage
	err := r.db.Get(&message, query, edit.NewMessage, edit.Id)
	if err != nil {
		if isDuplicateError(err) {
			return models.PvMessage{}, ErrDuplicateEntry
		}
		return models.PvMessage{}, ErrInternal
	}
	return message, nil
}

func (r *Repository) EditGroup(edit models.EditGroupMessage) (models.GroupMessage, error) {
	query := `UPDATE group_messages SET message = $1 WHERE id = $2 AND group_id = $3;`
	var message models.GroupMessage
	err := r.db.Get(&message, query, edit.NewMessage, edit.Id, edit.GroupId)
	if err != nil {
		if isDuplicateError(err) {
			return models.GroupMessage{}, ErrDuplicateEntry
		}
		return models.GroupMessage{}, ErrInternal
	}
	return message, nil
}

func isDuplicateError(err error) bool {

	return err != nil && (err.Error() == "pq: duplicate key value violates unique constraint")
}