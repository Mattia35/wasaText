package database
import (
	"database/sql"
)

func (db *appdbimpl) IsUserInConv(userId int, convId int) (bool, error) {
	var UserId int
	err := db.c.QueryRow(queryGetUsersConv, convId, userId).Scan(&UserId)
	if err != nil{
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}