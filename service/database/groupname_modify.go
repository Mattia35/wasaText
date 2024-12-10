package database
var query_CHANGEGROUPNAME = `UPDATE groupTable SET username = ? WHERE groupID = ?;`
func (db *appdbimpl) GroupnameModify (groupId int, groupname string) error{
	_, err := db.c.Exec(query_CHANGEGROUPNAME, groupname, groupId)
	return err
}