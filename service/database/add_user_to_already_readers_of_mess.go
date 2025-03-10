package database

var query_DELETEUSERTOREADERSOFMESS = `DELETE FROM checkMessTable WHERE messId = ? AND userId = ? AND convId = ?;`
func (db *appdbimpl) AddUserToListOfAlreadyReadersOfMess(messId int, userId int, convId int) error {
	_, err := db.c.Exec(query_DELETEUSERTOREADERSOFMESS, messId, userId, convId)
	return err
}