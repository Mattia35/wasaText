package database

var query_GETCONVBYUSERS = `
	SELECT COUNT(convp.ConversationId)
		FROM usersConvTable convuser, usersConvTable convus, convTable c
		WHERE convuser.ConversationId = convus.ConversationId
		  AND convuser.ConversationId = c.ConversationId
		  AND c.ConversationId = convus.ConversationId 
		  AND convuser.UserId = ? 
		  AND convus.UserId = ?  
		  AND c.GroupId IS NULL;`
func (db *appdbimpl) GetConvByUsers(userId int, destId int) (int, error) {
	var convNumber int
	err := db.c.QueryRow(query_GETCONVBYUSERS, userId, destId).Scan(&convNumber)
	if err != nil {
		return 0, err
	}
	return convNumber, nil
}