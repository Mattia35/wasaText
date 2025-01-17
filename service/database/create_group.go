package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"progetto.wasa/service/api/photoUtils"
)

var query_ADDGROUP = `INSERT INTO groupTable (groupId, username) VALUES (?, ?);`
var query_ADDGROUPUSER = `INSERT INTO usersGroupTable (groupId, userId) VALUES (?, ?);`
var query_ADDCONVUSER = `INSERT INTO usersConvTable (convId, userId) VALUES (?, ?);`
var query_MAXGROUPID = `SELECT MAX(groupId) FROM groupTable`

func (db *appdbimpl) CreateGroup(gr Group, convId int, userId int) (Group, error) {
	var group Group
	group.Username = gr.Username

	// ------FIND GROUPID---------//
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXGROUPID)
	if err != nil {
		return group, err
	}

	var maxID int
	row.Next()
	if row.Err() != nil {
		return group, err
	}

	err = row.Scan(&_maxID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return group, err
	}

	if !_maxID.Valid {
		maxID = 0
	} else {
		maxID = int(_maxID.Int64)
	}

	// --------SET GROUPID------------//
	group.GroupId = maxID + 1

	// --------CREATE GROUP FOLDER------------//
	path := "./storage/groups/" + fmt.Sprint(group.GroupId)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return group, err
	}

	// --------SET DEFAULT PROPIC------------//
	source, err := os.Open("./storage/default_profile_photo.jpg")
	if err != nil {
		return group, err
	}
	defer source.Close()

	destination, err := os.Create(photoUtils.GetGroupPhotoPath(group.GroupId))
	if err != nil {
		return group, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return group, err
	}

	// ------------INSERT GROUP--------------//
	_, err = db.c.Exec(query_ADDGROUP, group.GroupId, group.Username)
	if err != nil {
		return group, err
	}
	_, err = db.c.Exec(query_ADDGROUPUSER, group.GroupId, userId)
	if err != nil {
		return group, err
	}
	_, err = db.c.Exec(query_ADDCONVUSER, convId, userId)
	if err != nil {
		return group, err
	}
	return group, nil
}