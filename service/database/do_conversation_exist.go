package database

import (
	"database/sql"
)

var query_DOCONVERSATIONEXIST = `SELECT convId FROM convTable WHERE convId = ?;`

func (db *appdbimpl) DoConversationExist(ConvId int) (bool, error){
	var senderId int
	err := db.c.QueryRow(query_DOCONVERSATIONEXIST, ConvId).Scan(&senderId)
	if err != nil{
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}
