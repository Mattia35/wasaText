package database

var query_GETUSERBYNAME = `SELECT userId, username FROM User WHERE username = ?;`

func (database *appdbimpl) GetUserByName(username string) (User, error) {
	var user User
	err := database.c.QueryRow(query_GETUSERBYNAME, username).Scan(&user.UserId, &user.Username)
	return user, err
}