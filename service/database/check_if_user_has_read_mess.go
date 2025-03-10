package database
import (
	"database/sql"
)

var query_CHECKIFUSERHASREADMESS = `SELECT userId FROM checkMessTable WHERE messId = ? AND userId = ?;`
func (db *appdbimpl) CheckIfUserHasReadMess(messId int, userId int) (bool, error) {
	var UserId int
	err := db.c.QueryRow(query_CHECKIFUSERHASREADMESS, messId, userId).Scan(&UserId)
	if err != nil{
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, err
}