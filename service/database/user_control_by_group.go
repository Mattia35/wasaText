package database

import (
	"database/sql"
)

var query_USERCONTROLBYGROUPID = `SELECT groupId FROM usersGroupTable WHERE (groupId, userId) = (?,?);`

func (db *appdbimpl) UserControlByGroup(userId int, groupId int) (bool, error) {
	var GroupId int
	err := db.c.QueryRow(query_USERCONTROLBYGROUPID, groupId, userId).Scan(&GroupId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}
