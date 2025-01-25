package database

var query_USERCONTROLBYUSERNAME = `SELECT userId FROM userTable WHERE (username) = (?,?);`

func (database *appdbimpl) UserControlByUsername(username string) ( error) {
	_, err := database.c.Exec(query_USERCONTROLBYUSERNAME, username)
	return err
}