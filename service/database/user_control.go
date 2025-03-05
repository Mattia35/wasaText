package database

var query_USERCONTROLBYUSERNAME = `SELECT userId FROM userTable WHERE (username) = (?,?);`

func (db *appdbimpl) UserControlByUsername(username string) ( error) {
	_, err := db.c.Exec(query_USERCONTROLBYUSERNAME, username)
	return err
}