package database

var query_REMOVEGROUP = `DELETE FROM groupTable WHERE groupId = ?`

func (db *appdbimpl) RemoveGroup(groupId int) error {
	_, err := db.c.Exec(query_REMOVEGROUP, groupId)
	return err
}