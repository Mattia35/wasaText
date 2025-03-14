package database

var query_REMOVEUSERFROMCONV = `DELETE FROM usersConvTable WHERE userId = ? AND convId = ?;`

func (db *appdbimpl) RemoveUserFromConv(userId int, convId int) error {
	_, err := db.c.Exec(query_REMOVEUSERFROMCONV, userId, convId)
	return err
}
