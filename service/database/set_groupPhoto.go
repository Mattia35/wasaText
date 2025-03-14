package database

// update the status of a message
var query_UPDATEGROUPPHOTO = `UPDATE groupTable SET photo = ? WHERE groupId = ?;`

func (db *appdbimpl) SetGroupPhoto(groupId int, base64 string) error {
	_, err := db.c.Exec(query_UPDATEGROUPPHOTO, base64, groupId)
	return err
}
