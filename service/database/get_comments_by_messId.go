package database
import (
	"progetto.wasa/service/api/structions"
	
)

// get all comments of a message
var query_GETCOMMENTSBYMESSID = `SELECT commentId, messId, content, senderId, convId FROM commentTable WHERE messId = ?;`

func (db *appdbimpl) GetCommentsByMessId(messId int) ([]structions.Comment, error) {
	rows, err := db.c.Query(query_GETCOMMENTSBYMESSID, messId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []structions.Comment
	for rows.Next() {
		var comment structions.Comment
		err := rows.Scan(&comment.CommentId, &comment.MessageId, &comment.Content, &comment.SenderId, &comment.ConvId)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}