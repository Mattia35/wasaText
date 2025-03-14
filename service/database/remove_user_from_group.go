package database

// Query used to remove user from a group
var query_REMOVEUSERFROMGROUP = `DELETE FROM usersGroupTable WHERE groupId = ? AND userId = ?`

func (db *appdbimpl) RemoveUserFromGroup(userId int, groupId int) error {
	_, err := db.c.Exec(query_REMOVEUSERFROMGROUP, groupId, userId)
	return err
}
