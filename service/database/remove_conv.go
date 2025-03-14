package database

var query_REMOVECONV = `DELETE FROM convTable WHERE convId = ?`

func (db *appdbimpl) RemoveConv(convId int) error {
	_, err := db.c.Exec(query_REMOVECONV, convId)
	return err
}
