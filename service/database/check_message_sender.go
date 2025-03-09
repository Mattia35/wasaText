package database
var query_CHECKMESSAGESENDER = "SELECT senderId FROM messTable WHERE messId = ?"

func (db *appdbimpl) CheckMessageSender(messId int, userId int) (bool, error) {
	var senderId int
	err := db.c.QueryRow(query_CHECKMESSAGESENDER, messId).Scan(&senderId)
	if err != nil {
		return false, err
	}
	if senderId != userId {
		return false, err
	}
	return true, nil
}