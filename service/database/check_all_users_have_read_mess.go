package database

var query_ALLUSERSHAVEREADMESS = `SELECT userId FROM checkMessTable WHERE messId = ?;`
 func (db *appdbimpl) CheckAllUsersHaveReadMess(messId int) (bool, error) {
	var userIds []int
	rows, err := db.c.Query(query_ALLUSERSHAVEREADMESS, messId)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		if err != nil {
			return false, err
		}
		userIds = append(userIds, userId)
	}
	if len(userIds) == 0 {
		return true, nil
	}
	return false, nil
}