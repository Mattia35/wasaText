package database

import (
	"progetto.wasa/service/api/structions"
)

var query_GETGROUPBYGROUPID = `SELECT groupId, username, photo FROM GroupTable WHERE groupId = ?;`

func (db *appdbimpl) GetGroupByGroupId(groupId int) (structions.Group, error) {
	var group structions.Group
	err := db.c.QueryRow(query_GETGROUPBYGROUPID, groupId).Scan(&group.GroupId, &group.Username, &group.GroupPhoto)
	return group, err
}
