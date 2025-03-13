package database

import (
	"progetto.wasa/service/api/structions"
)
var query_GETMESSAGESBYCONVID = `SELECT messId, dateTime, IFNULL(text, ""), status, convId, IFNULL(photo, ""), IFNULL(gif, ""), senderId FROM messTable WHERE convId = ?`
func (db *appdbimpl) GetMessagesByConvId(convId int) ([]structions.Message, error) {
	rows, err := db.c.Query(query_GETMESSAGESBYCONVID, convId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []structions.Message
	for rows.Next() {
		var m structions.Message
		err := rows.Scan(&m.MessageId, &m.DateTime, &m.Text, &m.Status, &m.ConvId, &m.Photo, &m.Gif, &m.SenderId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}