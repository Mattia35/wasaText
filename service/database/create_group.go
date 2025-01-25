package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"progetto.wasa/service/api/photoUtils"
	"progetto.wasa/service/api/structions"
)

var query_ADDGROUP = `INSERT INTO groupTable (groupId, username) VALUES (?, ?);`
var query_ADDGROUPUSER = `INSERT INTO usersGroupTable (groupId, userId) VALUES (?, ?);`
var query_ADDCONVUSER = `INSERT INTO usersConvTable (convId, userId) VALUES (?, ?);`
var query_MAXGROUPID = `SELECT MAX(groupId) FROM groupTable`

func (db *appdbimpl) CreateGroup(gr structions.Group, convId int, userId int) (structions.Group, int, error) {
	var group structions.Group
	group.Username = gr.Username

	// ------FIND GROUPID---------//
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXGROUPID)
	if err != nil {
		return group, 0, err
	}

	var maxID int
	row.Next()
	if row.Err() != nil {
		return group, 0, err
	}

	err = row.Scan(&_maxID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return group, 0, err
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
		return group, maxID + 1, err
	}

	// --------SET DEFAULT PROPIC------------//
	source, err := os.Open("./storage/default_profile_photo.jpg")
	if err != nil {
		return group, maxID + 1, err
	}
	defer source.Close()

	destination, err := os.Create(photoUtils.GetGroupPhotoPath(group.GroupId))
	if err != nil {
		return group, maxID + 1, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return group, maxID + 1, err
	}

	// ------------INSERT GROUP--------------//
	_, err = db.c.Exec(query_ADDGROUP, group.GroupId, group.Username)
	if err != nil {
		return group, maxID + 1, err
	}
	_, err = db.c.Exec(query_ADDGROUPUSER, group.GroupId, userId)
	if err != nil {
		return group, maxID + 1, err
	}
	_, err = db.c.Exec(query_ADDCONVUSER, convId, userId)
	if err != nil {
		return group, maxID + 1, err
	}
	return group, maxID + 1, nil
}