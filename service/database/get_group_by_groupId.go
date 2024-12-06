package database

var query_GETGROUPBYGROUPID = `SELECT groupId, username FROM GroupTable WHERE groupId = ?;`

func (database *appdbimpl) GetGroupByGroupId(groupId int) (Group, error) {
	var group Group
	err := database.c.QueryRow(query_GETGROUPBYGROUPID, groupId).Scan(&group.GroupId, &group.Username)
	return group, err
}