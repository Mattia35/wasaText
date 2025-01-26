package database

import "progetto.wasa/service/api/structions"

var query_GETUSERBYID = `SELECT userId, username FROM UserTable WHERE userId = ?;`

func (database *appdbimpl) GetUserById(userId int) (structions.User, error) {
	var user structions.User
	err := database.c.QueryRow(query_GETUSERBYID, userId).Scan(&user.UserId, &user.Username)
	return user, err
}