package database

import (
	"progetto.wasa/service/api/structions"
)

// Query used to get user of a group
var query_GETCOMMENTINCONV = `SELECT commId, messId, content, senderId, convId FROM commentTable WHERE commId = ? AND messId = ? AND convId = ?`

func (db *appdbimpl) GetCommentById(commId int, messId int, convId int) (structions.Comment, error) {
	var comment structions.Comment
	err := db.c.QueryRow(query_GETCOMMENTINCONV, commId, messId, convId).Scan(&comment.CommentId, &comment.MessageId, &comment.Content, &comment.SenderId, &comment.ConvId)
	return comment, err
}
