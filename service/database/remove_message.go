package database

var query_REMOVEMESSAGE = "DELETE FROM messTable WHERE messId = ?"

func (db *appdbimpl) RemoveMessage(convId int) error {
	_, err := db.c.Exec(query_REMOVEMESSAGE, convId)
	if err != nil {
		return err
	}
	return err
}
