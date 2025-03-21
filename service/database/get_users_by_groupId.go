package database

import (
	"progetto.wasa/service/api/structions"
)

var query_GETUSERSBYGROUPID = `SELECT userId FROM usersGroupTable WHERE groupId = ?;`

func (db *appdbimpl) GetUsersByGroupId(groupId int) ([]structions.User, error) {
	var users []structions.User
	rows, err := db.c.Query(query_GETUSERSBYGROUPID, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user structions.User
		err = rows.Scan(&user.UserId)
		if err != nil {
			return nil, err
		}
		user, err = db.GetUserById(user.UserId)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
