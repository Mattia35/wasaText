package database
import (
	"progetto.wasa/service/api/structions"
	"database/sql"
	"errors"
)
var query_ADDCOMMENT = `INSERT INTO commentTable (commId, messId, content, senderId, convId) VALUES (?, ?, ?, ?, ?);`
var query_MAXCOMMID = `SELECT MAX(commId) FROM commentTable WHERE convId = ? AND messId = ?;`

func (db *appdbimpl) CreateComment(com structions.Comment) (structions.Comment, error) {
	var comment structions.Comment
	comment.CommentId = com.CommentId
	comment.ConvId = com.ConvId
	comment.SenderId = com.SenderId
	comment.Content = com.Content
	comment.MessageId = com.MessageId
	// ------FIND MESSID---------//
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXCOMMID, comment.ConvId, comment.MessageId)
	if err != nil {
		return comment, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return comment, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return comment, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	// --------SET USERID------------//
	comment.CommentId = maxID + 1
	// ------------INSERT USER--------------//
	_, err = db.c.Exec(query_ADDCOMMENT, comment.CommentId, comment.MessageId, comment.Content, comment.SenderId, comment.ConvId)
	if err != nil {
		return comment, err
	}

	return comment, nil
}