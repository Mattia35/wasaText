package database

import (
	"progetto.wasa/service/api/structions"
)

var query_USERCONTROLBYUSERNAME = `SELECT userId FROM userTable WHERE username = ?;`

func (db *appdbimpl) UserControlByUsername(username string) (structions.User, error) {
	var user structions.User
	err := db.c.QueryRow(query_USERCONTROLBYUSERNAME, username).Scan(&user.UserId)
	if err != nil {
		return user, err
	}
	return user, err
}
