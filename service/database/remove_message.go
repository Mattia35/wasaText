package database

var query_REMOVEMESSAGE = "DELETE FROM messTable WHERE messId = ? AND convId = ?"

func (db *appdbimpl) RemoveMessage(messId int, convId int) error {
	_, err := db.c.Exec(query_REMOVEMESSAGE, messId, convId)
	if err != nil {
		return err
	}
	return err
}
