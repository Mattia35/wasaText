package database

var query_CHANGEUSERNAME = `UPDATE userTable SET username = ? WHERE userID = ?;`

func (db *appdbimpl) UsernameModify(userId int, username string) error {
	_, err := db.c.Exec(query_CHANGEUSERNAME, username, userId)
	return err
}
