package database

var query_USERCONTROLBYGROUPID = `SELECT groupId FROM userGroupTable WHERE (groupId, userId) = (?,?);`

func (database *appdbimpl) UserControlByGroup(userId int, groupId int) (bool, error) {
	_, err := database.c.Exec(query_USERCONTROLBYGROUPID, groupId, userId)
	return true, err
}