package models

import (
	"github.com/Massad/gin-boilerplate/db"
)

//UserModel ...
type ConversationModel struct{}

type Conversation struct {
	ID           int64  `db:"id, primarykey, autoincrement" json:"id"`
	LastActivity int64  `db:"last_activity" json:"last_activity"`
	CreatedAt    int64  `db:"created_at" json:"created_at"`
	UpdatedAt    int64  `db:"updated_at" json:"updated_at"`
	UserID       string `db:"user_id" json:"user_id"`
	GroupID      string `db:"group_id" json:"group_id"`
}

type Conversations []Conversation

type ConversationDatasList struct {
	Data string `db:"data" json:"data"`
	Meta string `db:"meta" json:"meta"`
}

type ConversationDataStructure struct {
	CreatedAt  int64  `json:"created_at"`
	ID         int64  `json:"id"`
	Message    string `json:"message"`
	UpdatedAt  int64  `json:"updated_at"`
	SenderName string `json:"sender_name"`
	User       struct {
		Email string `json:"email"`
		ID    int64  `json:"id"`
		Name  string `json:"name"`
	} `json:"user"`
	UserID string `json:"user_id"`
}

type ConversationsDataStructure []ConversationDataStructure

func (m ConversationModel) InsertIntoConersations(userid string, groupuuid string) {

	_ = db.GetDB().QueryRow("INSERT INTO conversations(user_id, group_id) VALUES($1, $2)", userid, groupuuid)

}

func (m ConversationModel) GetConversations(userId string) (conversations Conversations) {

	_, _ = db.GetDB().Select(&conversations, "SELECT group_id FROM conversations WHERE user_id=$1 ORDER BY last_activity DESC", userId)

	return conversations

}

func (m ConversationModel) GetConversation(groupuuid string) (history string, err error) {

	var historySerialized []ConversationDatasList

	_, err = db.GetDB().Select(&historySerialized, "SELECT COALESCE(array_to_json(array_agg(row_to_json(d))), '[]') AS data, (SELECT row_to_json(n) FROM ( SELECT count(a.id) AS total FROM public.messages AS a WHERE a.group_id=$1 LIMIT 1 ) n ) AS meta FROM ( SELECT a.id, u.name AS sender_name, a.user_id, a.message, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.messages a LEFT JOIN public.user u ON a.user_id = u.user_id WHERE a.group_id=$1 ORDER BY a.id DESC) d", groupuuid)

	return historySerialized[0].Data, err
}

func (m ConversationModel) InsertIntoMessages(userid string, groupuuid string, message string) {

	_ = db.GetDB().QueryRow("INSERT INTO public.messages(user_id, group_id, message) VALUES($1, $2, $3)", userid, groupuuid, message)

}
