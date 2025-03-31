package database

import "progetto.wasa/service/api/structions"

// Query used to get user of a conversation
var queryGetOtherUsersConv = `SELECT userId FROM usersConvTable WHERE convId = ? AND userId != ?`

// Function
func (db *appdbimpl) GetOtherUserByConv(convId int, userId int) (structions.User, error) {
	var user structions.User
	err := db.c.QueryRow(queryGetOtherUsersConv, convId, userId).Scan(&user.UserId)
	return user, err
}
