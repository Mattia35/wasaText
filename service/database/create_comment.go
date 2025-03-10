package database
import (
	"progetto.wasa/service/api/structions"
)
var query_ADDCOMMENT = `INSERT INTO commentTable (messId, content, senderId, convId) VALUES (?, ?, ?, ?);`

func (db *appdbimpl) CreateComment(com structions.Comment) (structions.Comment, error) {
	var comment structions.Comment
	comment.ConvId = com.ConvId
	comment.SenderId = com.SenderId
	comment.Content = com.Content
	comment.MessageId = com.MessageId
	// ------------INSERT USER--------------//
	_, err := db.c.Exec(query_ADDCOMMENT, comment.MessageId, comment.Content, comment.SenderId, comment.ConvId)
	if err != nil {
		return comment, err
	}

	return comment, nil
}