package database

import (
	"progetto.wasa/service/api/structions"
)

var query_GETMESSBYID = `SELECT messId, dateTime, IFNULL(text, ""), status, convId, senderId, IFNULL(photo,""), IFNULL(gif,"") FROM messTable WHERE messId = ? AND ConvId = ?;`

func (db *appdbimpl) GetMessageById(LastMessage int, ConvId int) (structions.Message, error) {
	var mes structions.Message
	err := db.c.QueryRow(query_GETMESSBYID, LastMessage, ConvId).Scan(&mes.MessageId, &mes.DateTime, &mes.Text, &mes.Status, &mes.ConvId, &mes.SenderId, &mes.Photo, &mes.Gif)
	return mes, err
}
