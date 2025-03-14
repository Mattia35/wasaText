package database

var query_REMOVECOMMENT = "DELETE FROM commentTable WHERE commId = ? AND messId = ? AND convId = ?"

func (db *appdbimpl) RemoveComment(commId int, messId int, convId int) error {
	_, err := db.c.Exec(query_REMOVECOMMENT, commId, messId, convId)
	if err != nil {
		return err
	}
	return err
}
