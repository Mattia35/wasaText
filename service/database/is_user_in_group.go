package database
// Query used to get user of a group
var query_GETUSERGROUP = `SELECT userId FROM usersGroupTable WHERE groupId = ? AND userId = ?`

func (db *appdbimpl) IsUserInGroup(userId int, groupId int) (bool, error) {
	_,err := db.c.Exec(query_GETUSERGROUP, groupId, userId)
	if err != nil {
		return false, err
	}
	return true, err
}