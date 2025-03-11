package database

import (
	"database/sql"
	"errors"
	"encoding/base64"
	"io"
	"os"
	"progetto.wasa/service/api/structions"
)

var query_ADDGROUP = `INSERT INTO groupTable (groupId, username, photo) VALUES (?, ?, ?);`
var query_ADDGROUPUSER = `INSERT INTO usersGroupTable (groupId, userId) VALUES (?, ?);`
var query_ADDCONVUSER = `INSERT INTO usersConvTable (convId, userId) VALUES (?, ?);`
var query_MAXGROUPID = `SELECT MAX(groupId) FROM groupTable`

func (db *appdbimpl) CreateGroup(gr structions.Group, userId int) (structions.Group, error) {
	var group structions.Group
	group.Username = gr.Username

	// ------FIND GROUPID---------//
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXGROUPID)
	if err != nil {
		return group, err
	}

	var maxID int
	for row.Next(){
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
	}
	// --------SET GROUPID------------//
	group.GroupId = maxID + 1


	// Get the default profile photo
	file, err := os.Open("./storage/default_profile_photo.jpg")
	if err != nil {
		return group, err
	}
	defer file.Close()

	// Read the default profile photo
	data, err := io.ReadAll(file) 
		if err != nil {
			return group, err
		}
	
	// Encode the default profile photo
	group.GroupPhoto = base64.StdEncoding.EncodeToString(data)

	// ------------INSERT GROUP--------------//
	_, err = db.c.Exec(query_ADDGROUP, group.GroupId, group.Username, group.GroupPhoto)
	if err != nil {
		return group, err
	}
	_, err = db.c.Exec(query_ADDGROUPUSER, group.GroupId, userId)
	if err != nil {
		return group, err
	}
	return group, nil
}