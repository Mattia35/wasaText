package database


func (db *appdbimpl) AddUserToConv(userId int, convId int) (error) {
	_, err := db.c.Exec(query_ADDCONVUSER, convId, userId)
	return err
}