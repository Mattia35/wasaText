package database
// update the status of a message
var query_UPDATEMESSAGESTATUS = `UPDATE messTable SET status = ? WHERE messId = ?;`
func (db *appdbimpl) UpdateMessageStatus(messId int) error {
	_, err := db.c.Exec(query_UPDATEMESSAGESTATUS, true, messId)
	return err
}