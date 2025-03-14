package database

import "progetto.wasa/service/api/structions"

var QUERYGETCONVIDBYGROUPID = `SELECT convId, groupId, IFNULL(lastMessageId,0) FROM convTable WHERE groupId = ?;`

func (db *appdbimpl) GetConvByGroupId(groupId int) (structions.Conversation, error) {
	var conv structions.Conversation
	err := db.c.QueryRow(query_GETCONVBYID, groupId).Scan(&conv.ConvId, &conv.GroupId, &conv.LastMessage)
	return conv, err
}
