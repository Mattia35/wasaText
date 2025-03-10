package database
// update the status of a message
var query_UPDATECOMMENT = `UPDATE commentTable SET content = ? WHERE userId = ? AND convId = ? AND messId = ?;`
func (db *appdbimpl) UpdateComment(userId int, messId int, convId int, emoji string) error {
	_, err := db.c.Exec(query_UPDATECOMMENT, emoji, userId, convId, messId)
	return err
}