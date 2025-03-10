package database
import (
	"database/sql"
)
// Query used to get user of a group
var query_ISCOMMENTINCONV = `SELECT commId FROM commentTable WHERE commId = ? AND messId = ? AND convId = ?`

func (db *appdbimpl) IsCommentInConv(commId int, messId int, convId int) (bool, error) {
	var CommId int
	err := db.c.QueryRow(query_ISCOMMENTINCONV, commId, messId, convId).Scan(&CommId)
	if err != nil{
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}