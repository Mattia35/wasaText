package database

var query_USERCONTROLBYGROUPID = `SELECT groupId FROM userGroupTable WHERE (groupId, userId) = (?,?);`

func (db *appdbimpl) UserControlByGroup(userId int, groupId int) (bool, error) {
	_, err := db.c.Exec(query_USERCONTROLBYGROUPID, groupId, userId)
	return true, err
}