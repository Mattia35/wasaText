package database

import (
	"database/sql"
	"errors"
	"progetto.wasa/service/api/structions"
)

var query_ADDMESS = `INSERT INTO messTable (messId, dateTime, text, status, convId, Photo, senderId) VALUES (?, ?, ?, ?, ?, ?, ?);`
var query_MAXMESSID = `SELECT MAX(messId) FROM messTable WHERE convId = ?;`

func (db *appdbimpl) CreateMessage(mes structions.Message) (structions.Message, error) {
	var message structions.Message
	message.SenderId = mes.SenderId
	message.Text = mes.Text
	message.ConvId = mes.ConvId
	message.Photo = mes.Photo
	message.Status = mes.Status

	// ------FIND MESSID---------//
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXMESSID, message.ConvId)
	if err != nil {
		return message, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return message, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return message, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	// --------SET USERID------------//
	message.MessageId = maxID + 1


	// ------------INSERT USER--------------//
	_, err = db.c.Exec(query_ADDMESS, message.MessageId, message.DateTime, message.Text, message.Status, message.ConvId, message.Photo, message.SenderId)
	if err != nil {
		return message, err
	}

	return message, nil
}