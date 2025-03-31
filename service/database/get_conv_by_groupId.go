package database

import (
	"database/sql"

	"progetto.wasa/service/api/structions"
)

var QUERYGETCONVIDBYGROUPID = `SELECT convId, groupId, IFNULL(lastMessageId,0) FROM convTable WHERE groupId = ?;`

func (db *appdbimpl) GetConvByGroupId(groupId int) (bool, structions.Conversation, error) {
	var conv structions.Conversation
	err := db.c.QueryRow(QUERYGETCONVIDBYGROUPID, groupId).Scan(&conv.ConvId, &conv.GroupId, &conv.LastMessage)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, conv, nil
		}
		return false, conv, err
	}
	return true, conv, err
}
