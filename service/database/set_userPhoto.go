package database

// update the status of a message
var query_UPDATEUSERPHOTO = `UPDATE userTable SET photo = ? WHERE userId = ?;`
func (db *appdbimpl) SetUserPhoto(userId int, base64 string) error {
	_, err := db.c.Exec(query_UPDATEUSERPHOTO, base64, userId)
	return err
}