package database

import (
	"database/sql"
)

var queryCheckIfUserHasAlreadyCommented = `SELECT senderId FROM commentTable WHERE messId = ? AND senderId = ? AND convId = ?;`

func (db *appdbimpl) CheckIfUserHasAlreadyCommented(messId, userId, convId int) (bool, error) {
	var senderId int
	err := db.c.QueryRow(queryCheckIfUserHasAlreadyCommented, messId, userId, convId).Scan(&senderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}
