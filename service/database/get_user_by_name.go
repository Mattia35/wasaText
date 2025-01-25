package database

import "progetto.wasa/service/api/structions"

var query_GETUSERBYNAME = `SELECT userId, username FROM UserTable WHERE username = ?;`

func (database *appdbimpl) GetUserByName(username string) (structions.User, error) {
	var user structions.User
	err := database.c.QueryRow(query_GETUSERBYNAME, username).Scan(&user.UserId, &user.Username)
	return user, err
}