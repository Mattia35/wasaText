package database

import "progetto.wasa/service/api/structions"

var query_GETCONVBYID = `SELECT convId, IFNULL(groupId,0), IFNULL(lastMessageId,0) FROM convTable WHERE convId = ?;`

func (database *appdbimpl) GetConvById(convId int) (structions.Conversation, error) {
	var conv structions.Conversation
	err := database.c.QueryRow(query_GETCONVBYID, convId).Scan(&conv.ConvId, &conv.GroupId, &conv.LastMessage)
	return conv, err
}