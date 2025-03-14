package database

import "progetto.wasa/service/api/structions"

var query_GETUSERBYID = `SELECT userId, username, photo FROM UserTable WHERE userId = ?;`

func (db *appdbimpl) GetUserById(userId int) (structions.User, error) {
	var user structions.User
	err := db.c.QueryRow(query_GETUSERBYID, userId).Scan(&user.UserId, &user.Username, &user.UserPhoto)
	return user, err
}
