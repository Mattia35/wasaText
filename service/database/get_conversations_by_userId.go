package database

import (
	"progetto.wasa/service/api/structions"
)

var queryGetConversations = `SELECT convId FROM usersConvTable WHERE userId = ?`

func (db *appdbimpl) GetConversationsByUserId(userId int) ([]structions.Conversation, error) {
	var conversations []structions.Conversation
	rows, err := db.c.Query(queryGetConversations, userId)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var conv structions.Conversation
		err = rows.Scan(&conv.ConvId)
		if err != nil {
			return nil, err
		}
		conv, err = db.GetConvById(conv.ConvId)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}
