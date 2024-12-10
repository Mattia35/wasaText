package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"progetto.wasa/service/api/photoUtils"
)

var query_ADDUSER = `INSERT INTO userTable (userId, username) VALUES (?, ?);`
var query_MAXID = `SELECT MAX(userId) FROM userTable`

func (db *appdbimpl) CreateUser(u User) (User, error) {
	var user User
	user.Username = u.Username

	// ------FIND USERID---------//
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

	// --------SET USERID------------//
	user.UserId = maxID + 1

	// --------CREATE USER FOLDER------------//
	path := "./storage/" + fmt.Sprint(user.UserId) + "/conversations"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return user, err
	}

	// --------SET DEFAULT PROPIC------------//
	source, err := os.Open("./storage/default_profile_photo.jpg")
	if err != nil {
		return user, err
	}
	defer source.Close()

	destination, err := os.Create(photoUtils.GetUserPhotoPath(user.UserId))
	if err != nil {
		return user, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return user, err
	}

	// ------------INSERT USER--------------//
	_, err = db.c.Exec(query_ADDUSER, user.UserId, user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}