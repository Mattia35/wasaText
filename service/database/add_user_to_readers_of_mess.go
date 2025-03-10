package database

var query_ADDUSERTOREADERSOFMESS = `INSERT INTO checkMessTable (messId, userId, convId) VALUES (?, ?, ?);`
func (db *appdbimpl) AddUserToListOfReadersOfMess(messId int, userId int, convId int) error {
	_, err := db.c.Exec(query_ADDUSERTOREADERSOFMESS, messId, userId, convId)
	return err
}