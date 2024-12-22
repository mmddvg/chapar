package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mmddvg/chapar/pkg/errs"
	"mmddvg/chapar/pkg/models"
	"mmddvg/chapar/pkg/requests"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func isDuplicateError(err error) bool {

	return err != nil && (strings.Contains(err.Error(), "duplicate key value violates unique constraint"))
}

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo() *PostgresRepo {
	var str string
	Initialize(1) // TODO

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

func (r *PostgresRepo) SignUp(newUser requests.User) (models.User, error) {
	query := `INSERT INTO users (username, name, password) VALUES ($1, $2, $3) RETURNING id, username, name`
	var user models.User
	err := r.db.Get(&user, query, newUser.UserName, newUser.Name, newUser.Password)
	if err != nil {
		if isDuplicateError(err) {
			return models.User{}, errs.NewDuplicate("users", "")
		}
		return models.User{}, errs.NewUnexpected(err)
	}
	return user, nil
}

func (r *PostgresRepo) Get(userID uint64) (models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user models.User
	err := r.db.Get(&user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errs.NewNotFound("user", fmt.Sprint(userID))
		}
		return models.User{}, errs.NewUnexpected(err)
	}
	return user, nil
}

func (r *PostgresRepo) GetByUsername(userName string) (models.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	var user models.User
	err := r.db.Get(&user, query, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errs.NewNotFound("user", fmt.Sprint(userName))
		}
		return models.User{}, errs.NewUnexpected(err)
	}
	return user, nil
}
func (r *PostgresRepo) GetContacts(userID uint64) ([]models.Contact, error) {
	query := `SELECT 
    contacts.contact_id AS contact_id,
    users.username AS contact_username,
    private_chats.id AS pv_id
	FROM contacts JOIN 
    users ON contacts.contact_id = users.id
	LEFT JOIN 
    private_chats ON 
    (private_chats.user1 = contacts.user_id AND private_chats.user2 = contacts.contact_id)
    OR 
    (private_chats.user2 = contacts.user_id AND private_chats.user1 = contacts.contact_id)
	WHERE 
    contacts.user_id = $1;
`
	var contacts []models.Contact
	err := r.db.Select(&contacts, query, userID)
	if err != nil {
		return nil, errs.NewUnexpected(err)
	}
	return contacts, nil
}

func (r *PostgresRepo) IsContact(userID, contactID uint64) (bool, error) {
	query := `SELECT COUNT(*) > 0 FROM contacts WHERE user_id = $1 AND contact_id = $2`
	var exists bool
	err := r.db.Get(&exists, query, userID, contactID)
	if err != nil {
		return false, errs.NewUnexpected(err)
	}
	return exists, nil
}

func (r *PostgresRepo) AddContact(userID, contactID uint64) ([]models.Contact, error) {
	query := `INSERT INTO contacts (user_id, contact_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, userID, contactID)
	if err != nil {
		if isDuplicateError(err) {
			return nil, errs.NewDuplicate("contact", "")
		}
		return nil, errs.NewUnexpected(err)
	}
	return r.GetContacts(userID)
}

func (r *PostgresRepo) RemoveContact(userID, contactID uint64) ([]models.Contact, error) {
	query := `DELETE FROM contacts WHERE user_id = $1 AND contact_id = $2`
	_, err := r.db.Exec(query, userID, contactID)
	if err != nil {
		return nil, errs.NewUnexpected(err)
	}
	return r.GetContacts(userID)
}

func (r *PostgresRepo) AddGroupMember(groupId uint64, memberId uint64) (models.GroupMember, error) {
	var res models.GroupMember
	query := `INSERT INTO group_members(group_id,member_id) VALUES($1,$2) RETURNING *;`
	err := r.db.Get(&res, query, groupId, memberId)
	if err != nil {
		if isDuplicateError(err) {
			return res, errs.NewDuplicate("member", "id")
		}
		return res, errs.NewUnexpected(err)
	}
	return res, err
}

func (r *PostgresRepo) RemoveGroupMember(groupId uint64, memberId uint64) (models.GroupMember, error) {
	var res models.GroupMember
	query := `DELETE FROM group_members WHERE group_id = $1 AND member_id = $2 RETURNING *;`
	err := r.db.Get(&res, query, groupId, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errs.NewNotFound("member", fmt.Sprint(memberId))
		}
		return res, errs.NewUnexpected(err)
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
			return models.PrivateChat{}, errs.NewDuplicate("pv", "")
		}
		return models.PrivateChat{}, errs.NewUnexpected(err)
	}
	return chat, nil
}

// TODO : ppostgres hashmap index ?
func (r *PostgresRepo) GetPvOrCreate(user1Id, user2Id uint64) (models.PrivateChat, error) {
	query := `SELECT * FROM private_chats WHERE 
                (user1 = $1 AND user2 = $2) OR 
                (user1 = $2 AND user2 = $1);`

	var chat models.PrivateChat
	err := r.db.Get(&chat, query, user1Id, user2Id)
	if err != nil {
		if err == sql.ErrNoRows {
			insertQuery := `INSERT INTO private_chats (user1, user2) 
                            VALUES ($1, $2) 
                            RETURNING *;`
			err = r.db.Get(&chat, insertQuery, user1Id, user2Id)
			if err != nil {
				return chat, errs.NewUnexpected(err)
			}
		} else {
			return chat, errs.NewUnexpected(err)
		}
	}
	return chat, nil
}

func (r *PostgresRepo) WritePv(newMessage models.NewPvMessage) (models.PvMessage, error) {
	query := `INSERT INTO pv_messages VALUES ($1, $2, $3, $4, $5) RETURNING *;`
	var message models.PvMessage
	err := r.db.Get(&message, query, GenerateId(), newMessage.PvId, newMessage.SenderId, newMessage.Message, time.Now())
	if err != nil {
		return models.PvMessage{}, errs.NewUnexpected(err)
	}
	return message, nil
}

func (r *PostgresRepo) WriteGroup(newMessage models.NewGroupMessage) (models.GroupMessage, error) {
	query := `
		INSERT INTO group_messages 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, group_id, message, sender_id, created_at`
	var message models.GroupMessage
	err := r.db.Get(&message, query, GenerateId(), newMessage.GroupId, newMessage.Message, newMessage.SenderId, time.Now())
	if err != nil {
		return message, errs.NewUnexpected(err)
	}
	return message, nil
}

func (r *PostgresRepo) EditPv(edit models.EditPvMessage) (models.PvMessage, error) {
	query := `UPDATE pv_messages SET message = $1 WHERE id = $2 RETURNING *;`
	var message models.PvMessage
	err := r.db.Get(&message, query, edit.NewMessage, edit.Id)
	if err != nil {
		return models.PvMessage{}, errs.NewUnexpected(err)
	}
	return message, nil
}

func (r *PostgresRepo) EditGroup(edit models.EditGroupMessage) (models.GroupMessage, error) {
	query := `UPDATE group_messages SET message = $1 WHERE id = $2 AND group_id = $3;`
	var message models.GroupMessage
	err := r.db.Get(&message, query, edit.NewMessage, edit.Id, edit.GroupId)
	if err != nil {
		return models.GroupMessage{}, errs.NewUnexpected(err)
	}
	return message, nil
}

func (r *PostgresRepo) AddProfile(userId uint64, link string) ([]string, error) {
	var res []string
	query := "INSERT INTO user_profiles VALUES($1,$2,$3);"
	_, err := r.db.Exec(query, userId, link, time.Now())
	if err != nil {
		return res, errs.NewUnexpected(err)
	}

	err = r.db.Get(&res, "SELECT link FROM user_profiles WHERE user_id = $1;", userId)
	if err != nil {
		return res, errs.NewUnexpected(err)
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
		return res, errs.NewUnexpected(err)
	}

	err = r.db.Get(&res, "SELECT link FROM user_profiles WHERE user_id = $1 ORDER BY created_at;")
	if err != nil {
		return res, errs.NewUnexpected(err)
	}

	return res, nil
}

func (r *PostgresRepo) Block(userID, targetID uint64) (uint64, error) {
	query := `
		INSERT INTO blocked VALUES ($1, $2) RETURNING target_id;`
	err := r.db.Get(&targetID, query, userID, targetID)
	if err != nil {
		if isDuplicateError(err) {
			return 0, errs.NewDuplicate("blocked", fmt.Sprint(targetID))
		}
		return 0, errs.NewUnexpected(err)
	}
	return targetID, nil
}

func (r *PostgresRepo) UnBlock(userID, targetID uint64) (uint64, error) {
	query := `DELETE FROM blocked WHERE user_id = $1 AND target_id = $2;`
	_, err := r.db.Exec(query, userID, targetID)
	if err != nil {
		return 0, errs.NewUnexpected(err)
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
			return res, errs.NewDuplicate("group", "link")
		}
		return res, errs.NewUnexpected(err)
	}
	return res, nil
}

func (r *PostgresRepo) GetGroup(groupId uint64) (models.Group, error) {
	var group models.Group
	query := `SELECT * FROM groups WHERE id = $1`
	err := r.db.Get(&group, query, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Group{}, errs.NewNotFound("group", fmt.Sprint(groupId))
		}
		return models.Group{}, errs.NewUnexpected(err)
	}
	return group, nil
}

func (r *PostgresRepo) UpdateGroup(body requests.UpdateGroup) (models.Group, error) {
	var group models.Group
	query := `UPDATE groups SET title = $1 WHERE id = $3 RETURNING *`
	err := r.db.Get(&group, query, body.Name, body.GroupId)
	if err != nil {
		return models.Group{}, errs.NewUnexpected(err)
	}
	return group, nil
}

func (r *PostgresRepo) AddGroupProfile(groupId uint64, link string) (models.GroupProfile, error) {
	var profile models.GroupProfile
	query := `INSERT INTO group_profiles (g_id, link) VALUES ($1, $2) RETURNING *`
	err := r.db.Get(&profile, query, groupId, link)
	if err != nil {
		return models.GroupProfile{}, errs.NewUnexpected(err)
	}
	return profile, nil
}

func (r *PostgresRepo) RmGroupProfile(body requests.RmGroupProfile) (string, error) {
	var link string
	query := `DELETE FROM group_profiles
	WHERE (g_id, created_at) IN (
    SELECT g_id, created_at
    FROM group_profiles
    WHERE g_id = $1
    ORDER BY created_at
    OFFSET $2
) RETURNING link;`
	err := r.db.Get(&link, query, body.GroupId, body.NthCount)
	if err != nil {
		return "", errs.NewUnexpected(err)
	}
	return link, nil
}

func (r *PostgresRepo) GetGroupMembers(groupId uint64) ([]uint64, error) {
	var res []uint64
	query := "SELECT member_id FROM group_members WHERE group_id = $1;"
	err := r.db.Get(&res, query, groupId)
	if err != nil {
		return res, errs.NewUnexpected(err)
	}
	return res, err
}

func (r *PostgresRepo) IsMember(groupId, userId uint64) (bool, error) {
	var res uint64
	query := "SELECT member_id FROM group_members WHERE group_id = $1 AND member_id = $2;"

	err := r.db.Get(&res, query, groupId, userId)
	if err != nil {
		return false, errs.NewUnexpected(err)
	}

	return res == userId, nil
}

func (r *PostgresRepo) IsBlocked(userId uint64, targetId uint64) (bool, error) {
	var tmp uint64
	query := "SELECT user_id FROM blocked WHERE user_id = $1 AND target_id = $2;"

	err := r.db.Get(&tmp, query, userId, targetId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errs.NewUnexpected(err)
	}

	return tmp == userId, nil
}

func (r *PostgresRepo) SeenAck(messageId uint64) (models.PvMessage, error) {
	var res models.PvMessage
	query := "UPDATE pv_messages SET seen_at = NOW() WHERE id = $1 RETURNING *;"
	err := r.db.Get(&res, query, messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errs.NewNotFound("message", fmt.Sprint(messageId))
		}
		return res, errs.NewUnexpected(err)
	}
	return res, nil
}

func (r *PostgresRepo) GetChats(userId uint64) ([]models.PrivateChat, []models.Group, error) {
	users := []models.PrivateChat{}
	groups := []models.Group{}
	err := r.db.Select(&users, "SELECT * FROM private_chats WHERE user2 = $1 UNION ALL SELECT * FROM private_chats WHERE user1 = $1;", userId)
	if err != nil {
		return users, groups, errs.NewUnexpected(err)
	}

	err = r.db.Select(&groups, "SELECT * FROM groups WHERE id IN (SELECT group_id from group_members where member_id = $1);", userId)
	if err != nil {
		return users, groups, errs.NewUnexpected(err)
	}

	return users, groups, nil
}

func (r *PostgresRepo) GetPvMessages(id uint64) ([]models.PvMessage, error) {
	res := []models.PvMessage{}
	err := r.db.Select(&res, "SELECT * FROM pv_messages WHERE pv_id = $1;", id)
	if err != nil {
		return res, errs.NewUnexpected(err)
	}

	return res, err
}
func (r *PostgresRepo) GetGroupMessages(id uint64) ([]models.GroupMessage, error) {
	res := []models.GroupMessage{}
	err := r.db.Select(&res, "SELECT * FROM group_messages WHERE group_id = $1;", id)
	if err != nil {
		return res, errs.NewUnexpected(err)
	}

	return res, err
}
