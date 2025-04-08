package database

var query_CHECKMESSAGESENDER = "SELECT senderId FROM messTable WHERE messId = ? AND convId = ?"

func (db *appdbimpl) CheckMessageSender(messId int, userId int, convId int) (bool, error) {
	var senderId int
	err := db.c.QueryRow(query_CHECKMESSAGESENDER, messId, convId).Scan(&senderId)
	if err != nil {
		return false, err
	}
	if senderId != userId {
		return false, err
	}
	return true, nil
}
