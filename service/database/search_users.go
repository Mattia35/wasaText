package database

import "progetto.wasa/service/api/structions"

var query_SEARCHUSERS = `SELECT userId, username, photo FROM userTable WHERE username LIKE ?;`

func (db *appdbimpl) SearchUsers(query string) ([]structions.User, error) {
	rows, err := db.c.Query(query_SEARCHUSERS, query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structions.User
	for rows.Next() {
		var user structions.User
		err := rows.Scan(&user.UserId, &user.Username, &user.UserPhoto)
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
