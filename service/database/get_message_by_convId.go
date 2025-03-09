package database

import (
	"progetto.wasa/service/api/structions"
)
var query_GETMESSBYID = `SELECT messId, dateTime, text, status, convId, senderId, IFNULL(photo,"") FROM messTable WHERE messId = ? AND ConvId = ?;`
func (db *appdbimpl) GetMessageById(LastMessage int, ConvId int) (structions.Message, error){
	var mes structions.Message
	err := db.c.QueryRow(query_GETMESSBYID, LastMessage, ConvId).Scan(&mes.MessageId, &mes.DateTime, &mes.Text, &mes.Status, &mes.ConvId, &mes.SenderId, &mes.Photo)
	return mes, err
}