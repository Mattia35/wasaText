package database

import (
	"database/sql"
)

// Query used to get user of a group
var query_GETUSERGROUP = `SELECT userId FROM usersGroupTable WHERE groupId = ? AND userId = ?`

func (db *appdbimpl) IsUserInGroup(userId int, groupId int) (bool, error) {
	var UserId int
	err := db.c.QueryRow(query_GETUSERGROUP, groupId, userId).Scan(&UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}
