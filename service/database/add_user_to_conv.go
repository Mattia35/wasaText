package database


func (db *appdbimpl) AddUserToConv(userId int, convId int) (error) {
	_, err := db.c.Exec(query_ADDCONVUSER, convId, userId)
	if err != nil {
		return err
	}
	return err
}