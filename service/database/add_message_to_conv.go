package database

import "database/sql"

var query_ADDMESSCONV = `UPDATE convTable SET lastMessageId = ? WHERE convId = ?;`
func (db *appdbimpl) AddMessageToConv(MessageId int, ConvId int) (error) {
	isValid := MessageId != 0
	messId := sql.NullInt64{
		Int64: int64(MessageId), 
		Valid: isValid,
	}
	_, err := db.c.Exec(query_ADDMESSCONV, messId, ConvId)
	return err
}