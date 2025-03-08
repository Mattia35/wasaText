package database

import (
	"database/sql"
	"errors"

	"progetto.wasa/service/api/structions"
)

var query_ADDCONV = `INSERT INTO convTable (convId, groupId) VALUES (?,?);`
var query_ADDCONVNOGROUPID = `INSERT INTO convTable (convId) VALUES (?);`
var query_MAXCONVID = `SELECT MAX(convId) FROM convTable`

func (db *appdbimpl) CreateConversation(conv structions.Conversation) (structions.Conversation, error) {
	var conversation structions.Conversation
	conversation.LastMessage = conv.LastMessage

	// ------FIND CONVID---------//
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXCONVID)
	if err != nil {
		return conversation, err
	}

	var maxID int
	for row.Next(){
		if row.Err() != nil {
			return conversation, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return conversation, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}
	// --------SET CONVID------------//
	conversation.ConvId = maxID + 1
	conversation.GroupId = conv.GroupId
	if conversation.GroupId == 0 {
		_, err = db.c.Exec(query_ADDCONVNOGROUPID, conversation.ConvId)
		if err != nil {
			return conversation, err
		}
	} else{
		_, err = db.c.Exec(query_ADDCONV, conversation.ConvId, conversation.GroupId)
		if err != nil {
			return conversation, err
		}
	}
	return conversation, nil
}