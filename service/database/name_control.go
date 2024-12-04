package database

import (
	"database/sql"
	"errors"
)

var query_FINDUSERNAME = `SELECT username FROM User WHERE username = ?`

func (database *appdbimpl) NameControl(username string) (bool, error) {
	var existsName string
	err := database.c.QueryRow(query_FINDUSERNAME, username).Scan(&existsName)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return existsName != "", err
}