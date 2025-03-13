package database

var query_GETCONVERSATIONBYUSERS = `
	SELECT convuser.convId
		FROM usersConvTable convuser, usersConvTable convus, convTable c
		WHERE convuser.convId = convus.convId
		  AND convuser.ConvId = c.convId
		  AND c.convId = convus.convId 
		  AND convuser.userId = ? 
		  AND convus.userId = ?  
		  AND c.groupId IS NULL;`
func (db *appdbimpl) GetConversationByUsers(userId int, destId int) (int, error) {
	var convId int
	err := db.c.QueryRow(query_GETCONVERSATIONBYUSERS, userId, destId).Scan(&convId)
	if err != nil {
		return convId, err
	}
	return convId, err
}