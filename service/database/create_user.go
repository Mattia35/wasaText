package database

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"progetto.wasa/service/api/structions"
)

var query_ADDUSER = `INSERT INTO userTable (userId, username, photo) VALUES (?, ?, ?);`
var query_MAXID = `SELECT MAX(userId) FROM userTable`

func (db *appdbimpl) CreateUser(u structions.User) (structions.User, error) {
	var user structions.User
	user.Username = u.Username

	// Get the max ID
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXID)
	if err != nil {
		return user, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return user, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return user, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	// Set the new user ID
	user.UserId = maxID + 1

	// Get the default profile photo
	file, err := os.Open("./storage/default_profile_photo.jpg")
	if err != nil {
		return user, err
	}
	defer file.Close()

	// Read the default profile photo
	data, err := io.ReadAll(file)
	if err != nil {
		return user, err
	}

	// Encode the default profile photo
	user.UserPhoto = base64.StdEncoding.EncodeToString(data)

	// Add the user to the database
	_, err = db.c.Exec(query_ADDUSER, user.UserId, user.Username, user.UserPhoto)
	if err != nil {
		return user, err
	}

	return user, nil
}
