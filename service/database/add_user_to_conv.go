package database


func (db *appdbimpl) AddUserToConv(userId int, groupId int) (error) {
	_, err := db.c.Exec(query_ADDCONVUSER, groupId, userId)
	if err != nil {
		return err
	}
	return err
}