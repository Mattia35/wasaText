package database

import "progetto.wasa/service/api/structions"

var query_GETUSERBYNAME = `SELECT userId, username, photo FROM UserTable WHERE username = ?;`

func (db *appdbimpl) GetUserByName(username string) (structions.User, error) {
	var user structions.User
	err := db.c.QueryRow(query_GETUSERBYNAME, username).Scan(&user.UserId, &user.Username, &user.UserPhoto)
	return user, err
}