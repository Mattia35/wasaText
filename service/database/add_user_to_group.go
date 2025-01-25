package database


func (db *appdbimpl) AddUserToGroup(userId int, groupId int) (error) {
	_, err := db.c.Exec(query_ADDGROUPUSER, groupId, userId)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(query_ADDCONVUSER, groupId, userId)
	if err != nil {
		return err
	}
	return err
}