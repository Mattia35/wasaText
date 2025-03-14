package database

var query_GETCONVBYUSERS = `
	SELECT COUNT(convuser.convId)
		FROM usersConvTable convuser, usersConvTable convus, convTable c
		WHERE convuser.convId = convus.convId
		  AND convuser.ConvId = c.convId
		  AND c.convId = convus.convId 
		  AND convuser.userId = ? 
		  AND convus.userId = ?  
		  AND c.groupId IS NULL;`

func (db *appdbimpl) GetConvByUsers(userId int, destId int) (int, error) {
	var convNumber int
	err := db.c.QueryRow(query_GETCONVBYUSERS, userId, destId).Scan(&convNumber)
	if err != nil {
		return 0, err
	}
	return 1, err
}
