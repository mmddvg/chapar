package postgres

import (
	"database/sql"
	"errors"
	"log"
	"mmddvg/chapar/pkg/models"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func isDuplicateError(err error) bool {

	return err != nil && (err.Error() == "pq: duplicate key value violates unique constraint")
}

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo() *PostgresRepo {
	var str string

	if os.Getenv("POSTGRES_URI") != "" {
		str = os.Getenv("POSTGRES_URI")
	} else {
		str = "postgres://user:password@localhost:5432/chapar?sslmode=disable"
	}
	conn, err := sqlx.Connect("postgres", str)

	if err != nil {
		log.Fatal(err)
	}
	return &PostgresRepo{db: conn}
}

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrInternal       = errors.New("internal error")
)

func (r *PostgresRepo) SignUp(newUser models.NewUser) (models.User, error) {
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

func (r *PostgresRepo) Get(userID uint64) (models.User, error) {
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

func (r *PostgresRepo) GetContacts(userID uint64) ([]uint64, error) {
	query := `SELECT contact_id FROM contacts WHERE user_id = $1`
	var contacts []uint64
	err := r.db.Select(&contacts, query, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return contacts, nil
}

func (r *PostgresRepo) IsContact(userID, contactID uint64) (bool, error) {
	query := `SELECT COUNT(*) > 0 FROM contacts WHERE user_id = $1 AND contact_id = $2`
	var exists bool
	err := r.db.Get(&exists, query, userID, contactID)
	if err != nil {
		return false, ErrInternal
	}
	return exists, nil
}

func (r *PostgresRepo) AddContact(userID, contactID uint64) ([]uint64, error) {
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

func (r *PostgresRepo) RemoveContact(userID, contactID uint64) ([]uint64, error) {
	query := `DELETE FROM contacts WHERE user_id = $1 AND contact_id = $2`
	_, err := r.db.Exec(query, userID, contactID)
	if err != nil {
		return nil, ErrInternal
	}
	return r.GetContacts(userID)
}

func (r *PostgresRepo) AddGroupMember(groupId uint64, memberId uint64) (models.GroupMember, error) {
	var res models.GroupMember
	query := `INSERT INTO group_members(group_id,member_id) VALUES($1,$2) RETURNING *;`
	err := r.db.Get(&res, query, groupId, memberId)
	if err != nil {
		if isDuplicateError(err) {
			return res, ErrDuplicateEntry
		}
		return res, ErrInternal
	}
	return res, err
}

func (r *PostgresRepo) RemoveGroupMember(groupId uint64, memberId uint64) (models.GroupMember, error) {
	var res models.GroupMember
	query := `UPDATE group_members SET deleted_at = $1 WHERE group_id = $2 AND member_id = $3 RETURNING *;`
	err := r.db.Get(&res, query, time.Now(), groupId, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, ErrNotFound
		}
		return res, ErrInternal
	}
	return res, err
}

func (r *PostgresRepo) CreatePv(userID, targetID uint64) (models.PrivateChat, error) {
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

func (r *PostgresRepo) GetPv(targetID uint64) (models.PrivateChat, error) {
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

func (r *PostgresRepo) WritePv(newMessage models.NewPvMessage) (models.PvMessage, error) {
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

func (r *PostgresRepo) WriteGroup(newMessage models.NewGroupMessage) (models.GroupMessage, error) {
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

func (r *PostgresRepo) EditPv(edit models.EditPvMessage) (models.PvMessage, error) {
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

func (r *PostgresRepo) EditGroup(edit models.EditGroupMessage) (models.GroupMessage, error) {
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

func (r *PostgresRepo) AddProfile(userId uint64, link string) ([]string, error) {
	var res []string
	query := "INSERT INTO user_profiles VALUES($1,$2,$3);"
	_, err := r.db.Exec(query, userId, link, time.Now())
	if err != nil {
		return res, ErrInternal
	}

	err = r.db.Get(&res, "SELECT link FROM user_profiles WHERE user_id = $1;", userId)
	if err != nil {
		return res, ErrInternal
	}
	return res, nil
}
func (r *PostgresRepo) RemoveProfile(userId uint64, count uint) ([]string, error) {
	var res []string

	query := `
	DELETE FROM user_profiles
	WHERE (user_id, created_at) IN (
    SELECT user_id, created_at
    FROM user_profiles
    WHERE user_id = $1
    ORDER BY created_at
    OFFSET $2
);
`
	_, err := r.db.Exec(query, userId, count)
	if err != nil {
		return res, ErrInternal
	}

	err = r.db.Get(&res, "SELECT link FROM user_profiles WHERE user_id = $1 ORDER BY created_at;")
	if err != nil {
		return res, ErrInternal
	}

	return res, nil
}

func (r *PostgresRepo) Block(userID, targetID uint64) (uint64, error) {
	query := `
		INSERT INTO blocked VALUES ($1, $2) RETURNING target_id;`
	err := r.db.Get(&targetID, query, userID, targetID)
	if err != nil {
		if isDuplicateError(err) {
			return 0, ErrDuplicateEntry
		}
		return 0, ErrInternal
	}
	return targetID, nil
}

func (r *PostgresRepo) UnBlock(userID, targetID uint64) (uint64, error) {
	query := `DELETE FROM blocked WHERE user_id = $1 AND target_id = $2;`
	_, err := r.db.Exec(query, userID, targetID)
	if err != nil {
		return 0, ErrInternal
	}
	return targetID, nil
}

func (r *PostgresRepo) CreateGroup(userId uint64, title string, link string) (models.Group, error) {
	var res models.Group
	query := `
		INSERT INTO groups(title,link,owner_id) VALUES ($1, $2,$3) RETURNING *;`
	err := r.db.Get(&res, query, title, link, userId)
	if err != nil {
		if isDuplicateError(err) {
			return res, ErrDuplicateEntry
		}
		return res, ErrInternal
	}
	return res, nil
}
